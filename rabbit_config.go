package ocrworker

import (
	"flag"
)

type RabbitConfig struct {
	AmqpURI      string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Reliable     bool
}

func DefaultTestConfig() RabbitConfig {

	// Reliable: false due to major issues that would completely
	// wedge the rpc worker.  Setting the buffered channels length
	// higher would delay the problem, but then it would still happen later.

	rabbitConfig := RabbitConfig{
		AmqpURI:      "amqp://guest:guest@localhost:5672/",
		Exchange:     "open-ocr-exchange",
		ExchangeType: "direct",
		RoutingKey:   "decode-ocr",
		Reliable:     false, // setting to false because of observed issues
	}
	return rabbitConfig

}

type FlagFunction func()

func NoOpFlagFunction() FlagFunction {
	return func() {}
}

func DefaultConfigFlagsOverride(flagFunction FlagFunction) RabbitConfig {
	rabbitConfig := DefaultTestConfig()

	flagFunction()

	var AmqpURI string
	flag.StringVar(
		&AmqpURI,
		"amqp_uri",
		"",
		"The Amqp URI, eg: amqp://guest:guest@localhost:5672/",
	)

	flag.Parse()
	if len(AmqpURI) > 0 {
		rabbitConfig.AmqpURI = AmqpURI
	}

	return rabbitConfig

}
