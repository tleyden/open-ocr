package ocrworker

import (
	"encoding/json"
	"github.com/couchbaselabs/logg"
	"github.com/nu7hatch/gouuid"
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

type OcrRequest struct {
	ImgUrl     string
	EngineType OcrEngineType
}

func NewOcrRpcClient(rc RabbitConfig) (*OcrRpcClient, error) {
	ocrRpcClient := &OcrRpcClient{
		rabbitConfig: rc,
	}
	return ocrRpcClient, nil
}

func (c *OcrRpcClient) DecodeImageUrl(imgUrl string, eng OcrEngineType) (OcrResult, error) {

	var err error

	ocrRequest := OcrRequest{
		ImgUrl:     imgUrl,
		EngineType: eng,
	}
	ocrRequestJson, err := json.Marshal(ocrRequest)
	if err != nil {
		return OcrResult{}, err
	}
	logg.LogTo("OCR_CLIENT", "requestJson: %v", string(ocrRequestJson))

	correlationUuidRaw, err := uuid.NewV4()
	if err != nil {
		return OcrResult{}, err
	}
	correlationUuid := correlationUuidRaw.String()

	logg.LogTo("OCR_CLIENT", "dialing %q", c.rabbitConfig.AmqpURI)
	c.connection, err = amqp.Dial(c.rabbitConfig.AmqpURI)
	if err != nil {
		return OcrResult{}, err
	}

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

	err = c.subscribeCallbackQueue(correlationUuid)
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
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(ocrRequestJson),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			ReplyTo:         c.rabbitConfig.CallbackRoutingKey,
			CorrelationId:   correlationUuid,
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return OcrResult{}, nil
	}

	return OcrResult{}, nil
}

func (c OcrRpcClient) subscribeCallbackQueue(correlationUuid string) error {

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
		return err
	}

	// bind the callback queue to an exchange + routing key
	if err = c.channel.QueueBind(
		callbackQueue.Name,                // name of the queue
		c.rabbitConfig.CallbackRoutingKey, // bindingKey
		c.rabbitConfig.Exchange,           // sourceExchange
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	logg.LogTo("OCR_CLIENT", "callbackQueue name: %v", callbackQueue.Name)

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

	go c.handleRpcResponse(deliveries, correlationUuid)

	return nil

}

func (c OcrRpcClient) handleRpcResponse(deliveries <-chan amqp.Delivery, correlationUuid string) {
	logg.LogTo("OCR_CLIENT", "looping over deliveries..")
	for d := range deliveries {
		if d.CorrelationId == correlationUuid {
			logg.LogTo(
				"OCR_CLIENT",
				"got %dB delivery!!!!!!!: [%v] %q.  Reply to: %v",
				len(d.Body),
				d.DeliveryTag,
				d.Body,
				d.ReplyTo,
			)

			defer c.connection.Close()
			return

		} else {
			logg.LogTo("OCR_CLIENT", "ignoring delivery w/ correlation id: %v", d.CorrelationId)
		}

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
