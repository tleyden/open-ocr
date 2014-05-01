package ocrworker

import (
	"github.com/couchbaselabs/logg"
	"github.com/streadway/amqp"
)

type OcrRpcClient struct {
	rabbitConfig RabbitConfig
	connection   *amqp.Connection
	channel      *amqp.Channel
}

type OcrResult struct {
	Text string
}

func NewOcrRpcClient(rc RabbitConfig) (*OcrRpcClient, error) {
	ocrRpcClient := &OcrRpcClient{
		rabbitConfig: rc,
	}
	return ocrRpcClient, nil
}

func (c *OcrRpcClient) DecodeImageUrl(imgUrl string, eng OcrEngineType) (OcrResult, error) {

	var err error

	logg.LogTo("OCR_CLIENT", "dialing %q", c.rabbitConfig.AmqpURI)
	c.connection, err = amqp.Dial(c.rabbitConfig.AmqpURI)
	if err != nil {
		return OcrResult{}, err
	}

	// if we close the connection here, then we screw things up later
	// when subscribing to callback queue messages
	// defer c.connection.Close()

	c.channel, err = c.connection.Channel()
	if err != nil {
		return OcrResult{}, err
	}

	if err := c.channel.ExchangeDeclare(
		c.rabbitConfig.Exchange,     // name
		c.rabbitConfig.ExchangeType, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return OcrResult{}, err
	}

	// declare a callback queue where we will receive rpc responses
	callbackQueue, err := c.channel.QueueDeclare(
		c.rabbitConfig.CallbackQueueName, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return OcrResult{}, err
	}

	// bind the callback queue to an exchange + routing key
	if err = c.channel.QueueBind(
		callbackQueue.Name,                // name of the queue
		c.rabbitConfig.CallbackRoutingKey, // bindingKey
		c.rabbitConfig.Exchange,           // sourceExchange
		false, // noWait
		nil,   // arguments
	); err != nil {
		return OcrResult{}, err
	}

	// TODO: do we need to bind the callbackQueue to a key??

	logg.LogTo("OCR_CLIENT", "callbackQueue name: %v", callbackQueue.Name)

	// TODO: subscribe to this callback queue
	err = c.subscribeCallbackQueue(callbackQueue)
	if err != nil {
		return OcrResult{}, err
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if c.rabbitConfig.Reliable {
		if err := c.channel.Confirm(false); err != nil {
			return OcrResult{}, err
		}

		ack, nack := c.channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer confirmDelivery(ack, nack)
	}

	if err = c.channel.Publish(
		c.rabbitConfig.Exchange,   // publish to an exchange
		c.rabbitConfig.RoutingKey, // routing to 0 or more queues
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(imgUrl),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// ReplyTo:         callbackQueue.Name, Not working
			ReplyTo: c.rabbitConfig.CallbackRoutingKey,
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return OcrResult{}, nil
	}

	return OcrResult{}, nil
}

func (c OcrRpcClient) subscribeCallbackQueue(callbackQueue amqp.Queue) error {
	deliveries, err := c.channel.Consume(
		callbackQueue.Name, // name
		tag,                // consumerTag,
		true,               // noAck
		true,               // exclusive
		false,              // noLocal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	go c.handle(deliveries)

	return nil

}

func (c OcrRpcClient) handle(deliveries <-chan amqp.Delivery) {
	logg.LogTo("OCR_CLIENT", "looping over deliveries..")
	for d := range deliveries {
		logg.LogTo(
			"OCR_CLIENT",
			"got %dB delivery!!!!!!!: [%v] %q.  Reply to: %v",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
			d.ReplyTo,
		)

	}
}

func confirmDelivery(ack, nack chan uint64) {
	select {
	case tag := <-ack:
		logg.LogTo("OCR_CLIENT", "confirmed delivery, tag: %v", tag)
	case tag := <-nack:
		logg.LogTo("OCR_CLIENT", "failed to confirm delivery: %v", tag)
	}
}
