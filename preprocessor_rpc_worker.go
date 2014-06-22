package ocrworker

import (
	"encoding/json"
	"fmt"

	"github.com/couchbaselabs/logg"
	"github.com/streadway/amqp"
)

type PreprocessorRpcWorker struct {
	rabbitConfig RabbitConfig
	conn         *amqp.Connection
	channel      *amqp.Channel
	tag          string
	Done         chan error
}

const preprocessor_tag = "preprocessor" // TODO: should be unique for each worker instance (eg, uuid)

func NewPreprocessorRpcWorker(rc RabbitConfig) (*PreprocessorRpcWorker, error) {
	preprocessorRpcWorker := &PreprocessorRpcWorker{
		rabbitConfig: rc,
		conn:         nil,
		channel:      nil,
		tag:          preprocessor_tag,
		Done:         make(chan error),
	}
	return preprocessorRpcWorker, nil
}

func (w PreprocessorRpcWorker) Run() error {

	var err error

	logg.LogTo("PREPROCESSOR_WORKER", "Run() called...")

	logg.LogTo("PREPROCESSOR_WORKER", "dialing %q", w.rabbitConfig.AmqpURI)
	w.conn, err = amqp.Dial(w.rabbitConfig.AmqpURI)
	if err != nil {
		return err
	}

	go func() {
		fmt.Printf("closing: %s", <-w.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	logg.LogTo("PREPROCESSOR_WORKER", "got Connection, getting Channel")
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

	bindingKey := "stroke-width-transform"

	if err = w.channel.QueueBind(
		queue.Name,              // name of the queue
		bindingKey,              // bindingKey
		w.rabbitConfig.Exchange, // sourceExchange
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	logg.LogTo("PREPROCESSOR_WORKER", "Queue bound to Exchange, starting Consume (consumer tag %q, binding key: %v)", preprocessor_tag, bindingKey)
	deliveries, err := w.channel.Consume(
		queue.Name,       // name
		preprocessor_tag, // consumerTag,
		true,             // noAck
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	go w.handle(deliveries, w.Done)

	return nil
}

func (w *PreprocessorRpcWorker) Shutdown() error {
	// will close() the deliveries channel
	if err := w.channel.Cancel(w.tag, true); err != nil {
		return fmt.Errorf("Worker cancel failed: %s", err)
	}

	if err := w.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logg.LogTo("PREPROCESSOR_WORKER", "Shutdown OK")

	// wait for handle() to exit
	return <-w.Done
}

func (w *PreprocessorRpcWorker) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		logg.LogTo(
			"PREPROCESSOR_WORKER",
			"got %d byte delivery: [%v].  Reply to: %v",
			len(d.Body),
			d.DeliveryTag,
			d.ReplyTo,
		)

		err := w.handleDelivery(d)
		if err != nil {
			msg := "Error handling delivery in preprocessor.  Error: %v"
			logg.LogError(fmt.Errorf(msg, err))
		}

	}
	logg.LogTo("PREPROCESSOR_WORKER", "handle: deliveries channel closed")
	done <- fmt.Errorf("handle: deliveries channel closed")
}

func (w *PreprocessorRpcWorker) handleDelivery(d amqp.Delivery) error {

	ocrRequest := OcrRequest{}
	err := json.Unmarshal(d.Body, &ocrRequest)
	if err != nil {
		msg := "Error unmarshaling json: %v.  Error: %v"
		errMsg := fmt.Sprintf(msg, string(d.Body), err)
		logg.LogError(fmt.Errorf(errMsg))
		return err
	}

	logg.LogTo("PREPROCESSOR_WORKER", "ocrRequest before: %v", ocrRequest)
	routingKey := ocrRequest.nextPreprocessor(w.rabbitConfig.RoutingKey)
	logg.LogTo("PREPROCESSOR_WORKER", "publishing with routing key %q", routingKey)
	logg.LogTo("PREPROCESSOR_WORKER", "ocrRequest after: %v", ocrRequest)

	// TODO: process the image and then re-marshal ocrRequest

	ocrRequestJson, err := json.Marshal(ocrRequest)
	if err != nil {
		return err
	}

	logg.LogTo("PREPROCESSOR_WORKER", "sendRpcResponse to: %v", routingKey)

	if err := w.channel.Publish(
		w.rabbitConfig.Exchange, // publish to an exchange
		routingKey,              // routing to 0 or more queues
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(ocrRequestJson),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			ReplyTo:         d.ReplyTo,
			CorrelationId:   d.CorrelationId,
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return err
	}
	logg.LogTo("PREPROCESSOR_WORKER", "handleDelivery succeeded")

	return nil
}
