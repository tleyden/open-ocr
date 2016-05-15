package ocrworker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/couchbaselabs/logg"
	"github.com/streadway/amqp"
)

type OcrRpcWorker struct {
	rabbitConfig RabbitConfig
	conn         *amqp.Connection
	channel      *amqp.Channel
	tag          string
	Done         chan error
}

const tag = "foo" // TODO: should be unique for each worker instance (eg, uuid)

func NewOcrRpcWorker(rc RabbitConfig) (*OcrRpcWorker, error) {
	ocrRpcWorker := &OcrRpcWorker{
		rabbitConfig: rc,
		conn:         nil,
		channel:      nil,
		tag:          tag,
		Done:         make(chan error),
	}
	return ocrRpcWorker, nil
}

func (w OcrRpcWorker) Run() error {

	var err error

	logg.LogTo("OCR_WORKER", "Run() called...")

	logg.LogTo("OCR_WORKER", "dialing %q", w.rabbitConfig.AmqpURI)
	w.conn, err = amqp.Dial(w.rabbitConfig.AmqpURI)
	if err != nil {
		logg.LogTo("OCR_WORKER", "error connecting to rabbitmq %v", err)
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

	// just use the routing key as the queue name, since there's no reason
	// to have a different name
	queueName := w.rabbitConfig.RoutingKey

	queue, err := w.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	logg.LogTo("OCR_WORKER", "binding to: %v", w.rabbitConfig.RoutingKey)

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

	go w.handle(deliveries, w.Done)

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
	return <-w.Done
}

func (w *OcrRpcWorker) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		logg.LogTo(
			"OCR_WORKER",
			"got %d byte delivery: [%v]. Routing key: %v  Reply to: %v",
			len(d.Body),
			d.DeliveryTag,
			d.RoutingKey,
			d.ReplyTo,
		)

		ocrResult, err := w.resultForDelivery(d)
		if err != nil {
			msg := "Error generating ocr result.  Error: %v"
			logg.LogError(fmt.Errorf(msg, err))
		}

		logg.LogTo("OCR_WORKER", "Sending rpc response: %v", ocrResult)
		err = w.sendRpcResponse(ocrResult, d.ReplyTo, d.CorrelationId)
		if err != nil {
			msg := "Error returning ocr result: %v.  Error: %v"
			logg.LogError(fmt.Errorf(msg, ocrResult, err))
			// if we can't send our response, let's just abort
			done <- err
			break
		}

	}
	logg.LogTo("OCR_WORKER", "handle: deliveries channel closed")
	done <- fmt.Errorf("handle: deliveries channel closed")
}

func (w *OcrRpcWorker) resultForDelivery(d amqp.Delivery) (OcrResult, error) {

	ocrRequest := OcrRequest{}
	ocrResult := OcrResult{Text: "Error"}
	err := json.Unmarshal(d.Body, &ocrRequest)
	if err != nil {
		msg := "Error unmarshaling json: %v.  Error: %v"
		errMsg := fmt.Sprintf(msg, string(d.Body), err)
		logg.LogError(fmt.Errorf(errMsg))
		ocrResult.Text = errMsg
		return ocrResult, err
	}

	ocrEngine := NewOcrEngine(ocrRequest.EngineType)

	ocrResult, err = ocrEngine.ProcessRequest(ocrRequest)

	if err != nil {
		msg := "Error processing image url: %v.  Error: %v"
		errMsg := fmt.Sprintf(msg, ocrRequest.ImgUrl, err)
		logg.LogError(fmt.Errorf(errMsg))
		ocrResult.Text = errMsg
		return ocrResult, err
	}

	return ocrResult, nil

}

func (w *OcrRpcWorker) sendRpcResponse(r OcrResult, replyTo string, correlationId string) error {

	if w.rabbitConfig.Reliable {
		// Do not use w.rabbitConfig.Reliable=true due to major issues
		// that will completely  wedge the rpc worker.  Setting the
		// buffered channels length higher would delay the problem,
		// but then it would still happen later.
		if err := w.channel.Confirm(false); err != nil {
			return err
		}

		ack, nack := w.channel.NotifyConfirm(make(chan uint64, 100), make(chan uint64, 100))

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
	logg.LogTo("OCR_WORKER", "awaiting delivery confirmation ...")
	select {
	case tag := <-ack:
		logg.LogTo("OCR_WORKER", "confirmed delivery, tag: %v", tag)
	case tag := <-nack:
		logg.LogTo("OCR_WORKER", "failed to confirm delivery: %v", tag)
	case <-time.After(RPC_RESPONSE_TIMEOUT):
		// this is bad, the worker will probably be dsyfunctional
		// at this point, so panic
		logg.LogPanic("timeout trying to confirm delivery")
	}
}
