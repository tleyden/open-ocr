package ocrworker

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	"github.com/streadway/amqp"
)

type OcrRpcWorker struct {
	rabbitConfig RabbitConfig
	conn         *amqp.Connection
	channel      *amqp.Channel
	tag          string
	done         chan error
}

const tag = "foo" // TODO: should be unique for each worker instance (eg, uuid)

func NewOcrRpcWorker(rc RabbitConfig) (*OcrRpcWorker, error) {
	ocrRpcWorker := &OcrRpcWorker{
		rabbitConfig: rc,
		conn:         nil,
		channel:      nil,
		tag:          tag,
		done:         make(chan error),
	}
	return ocrRpcWorker, nil
}

func (w OcrRpcWorker) Run() error {

	var err error

	logg.LogTo("OCR_WORKER", "Run() called")

	logg.LogTo("OCR_WORKER", "dialing %q", w.rabbitConfig.AmqpURI)
	w.conn, err = amqp.Dial(w.rabbitConfig.AmqpURI)
	if err != nil {
		return err
	}

	go func() {
		fmt.Printf("closing: %s", <-w.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	logg.LogTo("OCR_WORKER", "got Connection, getting Channel")
	w.channel, err = w.conn.Channel()
	if err != nil {
		return err
	}

	if err = w.channel.ExchangeDeclare(
		w.rabbitConfig.Exchange,     // name of the exchange
		w.rabbitConfig.ExchangeType, // type
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	queue, err := w.channel.QueueDeclare(
		w.rabbitConfig.QueueName, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	if err = w.channel.QueueBind(
		queue.Name,                // name of the queue
		w.rabbitConfig.RoutingKey, // bindingKey
		w.rabbitConfig.Exchange,   // sourceExchange
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	logg.LogTo("OCR_WORKER", "Queue bound to Exchange, starting Consume (consumer tag %q)", tag)
	deliveries, err := w.channel.Consume(
		queue.Name, // name
		tag,        // consumerTag,
		true,       // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	go w.handle(deliveries, w.done)

	return nil
}

func (w *OcrRpcWorker) Shutdown() error {
	// will close() the deliveries channel
	if err := w.channel.Cancel(w.tag, true); err != nil {
		return fmt.Errorf("Worker cancel failed: %s", err)
	}

	if err := w.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logg.LogTo("OCR_WORKER", "Shutdown OK")

	// wait for handle() to exit
	return <-w.done
}

func (w *OcrRpcWorker) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		logg.LogTo(
			"OCR_WORKER",
			"got %dB delivery: [%v] %q.  Reply to: %v",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
			d.ReplyTo,
		)

		engineType := ENGINE_TESSERACT // TODO: the engine type should be specified in the message
		ocrEngine := NewOcrEngine(engineType)
		imgUrl := string(d.Body)
		ocrResult, err := ocrEngine.ProcessImageUrl(imgUrl)
		if err != nil {
			msg := "Error processing image url: %v.  Error: %v"
			logg.LogError(fmt.Errorf(msg, imgUrl, err))
			done <- err
			break
		}

		err = w.sendRpcResponse(ocrResult, d.ReplyTo, d.CorrelationId)
		if err != nil {
			msg := "Error returning ocr result: %v.  Error: %v"
			logg.LogError(fmt.Errorf(msg, ocrResult, err))
			done <- err
			break
		}

	}
	logg.LogTo("OCR_WORKER", "handle: deliveries channel closed")
	done <- nil
}

func (w *OcrRpcWorker) sendRpcResponse(r OcrResult, replyTo string, correlationId string) error {
	if w.rabbitConfig.Reliable {
		if err := w.channel.Confirm(false); err != nil {
			return err
		}

		ack, nack := w.channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer confirmDeliveryWorker(ack, nack)
	}

	logg.LogTo("OCR_WORKER", "sendRpcResponse to: %v", replyTo)
	if err := w.channel.Publish(
		w.rabbitConfig.Exchange, // publish to an exchange
		replyTo,                 // routing to 0 or more queues
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(r.Text),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			CorrelationId:   correlationId,
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return err
	}
	logg.LogTo("OCR_WORKER", "sendRpcResponse succeeded")
	return nil

}

func confirmDeliveryWorker(ack, nack chan uint64) {
	select {
	case tag := <-ack:
		logg.LogTo("OCR_WORKER", "confirmed delivery, tag: %v", tag)
	case tag := <-nack:
		logg.LogTo("OCR_WORKER", "failed to confirm delivery: %v", tag)
	}
}
