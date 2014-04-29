package ocrworker

type RabbitConfig struct {
	AmqpURI      string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Reliable     bool
}
