package ocrworker

import (
	"github.com/couchbaselabs/logg"
	"github.com/streadway/amqp"
)

type OcrRpcClient struct {
	rabbitConfig RabbitConfig
}

type OcrResult struct {
}

func NewOcrRpcClient(rc RabbitConfig) (*OcrRpcClient, error) {
	ocrRpcClient := &OcrRpcClient{
		rabbitConfig: rc,
	}
	return ocrRpcClient, nil
}

func (c OcrRpcClient) DecodeImageUrl(imgUrl string, eng OcrEngineType) (OcrResult, error) {

	logg.LogTo("OCR_RPC", "dialing %q", c.rabbitConfig.AmqpURI)
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
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return OcrResult{}, err
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if reliable {
		if err := channel.Confirm(false); err != nil {
			return OcrResult{}, err
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer confirmOne(ack, nack)
	}

	if err = channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return OcrResult{}, nil
	}

	return OcrResult{}, nil
}
