package ocrworker

import (
	"testing"

	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
)

func init() {
	logg.LogKeys["TEST"] = true
	logg.LogKeys["OCR"] = true
	logg.LogKeys["OCR_CLIENT"] = true
	logg.LogKeys["OCR_WORKER"] = true
	logg.LogKeys["PREPROCESSOR_WORKER"] = true
	logg.LogKeys["OCR_HTTP"] = true
	logg.LogKeys["OCR_TESSERACT"] = true
}

func rabbitConfigForTests() RabbitConfig {
	rabbitConfig := DefaultTestConfig()
	return rabbitConfig
}

// This test assumes that rabbit mq is running
func DisabledTestOcrRpcClientIntegration(t *testing.T) {

	// TODO: serve this up through a fake webserver
	// that reads from the filesystem
	testImageUrl := "http://localhost:8080/img"

	rabbitConfig := rabbitConfigForTests()

	// kick off a worker
	// this would normally happen on a different machine ..
	ocrWorker, err := NewOcrRpcWorker(rabbitConfig)
	if err != nil {
		logg.LogTo("TEST", "err: %v", err)
	}
	ocrWorker.Run()

	ocrClient, err := NewOcrRpcClient(rabbitConfig)
	if err != nil {
		logg.LogTo("TEST", "err: %v", err)
	}
	assert.True(t, err == nil)

	for i := 0; i < 50; i++ {

		ocrRequest := OcrRequest{ImgUrl: testImageUrl, EngineType: ENGINE_MOCK}
		decodeResult, err := ocrClient.DecodeImage(ocrRequest)
		if err != nil {
			logg.LogTo("TEST", "err: %v", err)
		}
		assert.True(t, err == nil)
		logg.LogTo("TEST", "decodeResult: %v", decodeResult)
		assert.Equals(t, decodeResult.Text, MOCK_ENGINE_RESPONSE)

	}

}
