package ocrworker

type RabbitConfig struct {
	AmqpURI            string
	Exchange           string
	ExchangeType       string
	RoutingKey         string
	CallbackRoutingKey string
	Reliable           bool
	QueueName          string
	CallbackQueueName  string
}
