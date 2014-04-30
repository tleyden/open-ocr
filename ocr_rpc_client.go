package ocrworker

import (
	"github.com/couchbaselabs/logg"
	"github.com/streadway/amqp"
)

type OcrRpcClient struct {
	rabbitConfig RabbitConfig
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

func (c OcrRpcClient) DecodeImageUrl(imgUrl string, eng OcrEngineType) (OcrResult, error) {

	logg.LogTo("OCR_CLIENT", "dialing %q", c.rabbitConfig.AmqpURI)
	connection, err := amqp.Dial(c.rabbitConfig.AmqpURI)
	if err != nil {
		return OcrResult{}, err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return OcrResult{}, err
	}

	if err := channel.ExchangeDeclare(
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
	callbackQueue, err := channel.QueueDeclare(
		c.rabbitConfig.CallbackQueueName, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return OcrResult{}, err
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if c.rabbitConfig.Reliable {
		if err := channel.Confirm(false); err != nil {
			return OcrResult{}, err
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer confirmDelivery(ack, nack)
	}

	if err = channel.Publish(
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
			ReplyTo:         callbackQueue.Name,
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return OcrResult{}, nil
	}

	return OcrResult{}, nil
}

func confirmDelivery(ack, nack chan uint64) {
	select {
	case tag := <-ack:
		logg.LogTo("OCR_CLIENT", "confirmed delivery, tag: %v", tag)
	case tag := <-nack:
		logg.LogTo("OCR_CLIENT", "failed to confirm delivery: %v", tag)
	}
}
