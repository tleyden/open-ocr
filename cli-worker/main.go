package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/ocr-worker"
)

// This assumes that there is a rabbit mq running
// To test it, fire up a webserver and send it a curl request

func init() {
	logg.LogKeys["OCR_CLIENT"] = true
	logg.LogKeys["OCR_WORKER"] = true
	logg.LogKeys["OCR_HTTP"] = true
	logg.LogKeys["OCR_TESSERACT"] = true
}

func main() {

	rabbitConfig := ocrworker.RabbitConfig{
		AmqpURI:            "amqp://guest:guest@localhost:5672/",
		Exchange:           "test-exchange",
		ExchangeType:       "direct",
		RoutingKey:         "test-key",
		CallbackRoutingKey: "callback-key",
		Reliable:           true,
		QueueName:          "test-queue",
		CallbackQueueName:  "callback-queue",
	}

	ocrWorker, err := ocrworker.NewOcrRpcWorker(rabbitConfig)
	if err != nil {
		logg.LogPanic("Could not create rpc worker")
	}
	ocrWorker.Run()

	select {}

	// blockForeverChan := make(chan bool)
	// <-blockForeverChan

}
