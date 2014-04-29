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

	rabbitConfig := RabbitConfig{}
	ocrClient, err := NewOcrRpcClient(rabbitConfig)
	assert.True(t, err == nil)
	decodeResult, err := ocrClient.DecodeImageUrl("http://foo.png", ENGINE_TESSERACT)
	assert.True(t, err == nil)
	logg.LogTo("TEST", "decodeResult: %v", decodeResult)

	// TODO: add assertions on decodeResult ..

}
