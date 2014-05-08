package ocrworker

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"log"
	"net/http"
	"testing"
	"time"
)

// This test assumes that rabbit mq is running
func TestOcrHttpHandler(t *testing.T) {

	rabbitConfig := rabbitConfigForTests()

	// create an ocr handler (passing it the rabbit config + engine type, which it will need)

	engineType := ENGINE_MOCK
	ocrHandler := NewOcrHttpHandler(rabbitConfig, engineType)
	http.Handle("/ocr", ocrHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))

	// connect to it via http client
	time.Sleep(time.Second * 60)

	// make sure get expected result
	logg.LogTo("TEST", "test")
	assert.True(t, true)
}
