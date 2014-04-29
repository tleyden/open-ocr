package ocrworker

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"testing"
)

func init() {
	logg.LogKeys["TEST"] = true
}

func TestOcrRpcClientIntegration(t *testing.T) {

	// assumes that rabbit mq is running

	rabbitConfig := RabbitConfig{
		AmqpURI:      "amqp://guest:guest@localhost:5672/",
		Exchange:     "test-exchange",
		ExchangeType: "direct",
		RoutingKey:   "test-key",
		Reliable:     true,
	}

	ocrClient, err := NewOcrRpcClient(rabbitConfig)
	if err != nil {
		logg.LogTo("TEST", "err: %v", err)
	}
	assert.True(t, err == nil)
	decodeResult, err := ocrClient.DecodeImageUrl("http://foo.png", ENGINE_TESSERACT)
	if err != nil {
		logg.LogTo("TEST", "err: %v", err)
	}
	assert.True(t, err == nil)
	logg.LogTo("TEST", "decodeResult: %v", decodeResult)

	// TODO: add assertions on decodeResult ..

}
