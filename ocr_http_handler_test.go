package ocrworker

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
)

// This test assumes that rabbit mq is running
func DisabledTestOcrHttpHandlerIntegration(t *testing.T) {

	rabbitConfig := rabbitConfigForTests()

	err := spawnOcrWorker(rabbitConfig)
	if err != nil {
		logg.LogPanic("Could not spawn ocr worker")
	}

	// add a handler to serve up an image from the filesystem.
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "refactoring.png")
	})

	http.Handle("/ocr", NewOcrHttpHandler(rabbitConfig))

	go http.ListenAndServe(":8081", nil)

	logg.LogTo("TEST", "test1")

	ocrRequest := OcrRequest{
		ImgUrl:     "http://localhost:8081/img",
		EngineType: ENGINE_MOCK,
	}
	jsonBytes, err := json.Marshal(ocrRequest)
	if err != nil {
		logg.LogPanic("Could not marshal OcrRequest")
	}

	reader := bytes.NewReader(jsonBytes)

	resp, err := http.Post("http://localhost:8081/ocr", "application/json", reader)
	assert.True(t, err == nil)
	logg.LogTo("TEST", "resp: %v", resp)

	// connect to it via http client
	logg.LogTo("TEST", "Sleep for 60s")
	time.Sleep(time.Second * 60)

	// make sure get expected result

	assert.True(t, true)
}

func spawnOcrWorker(rabbitConfig RabbitConfig) error {

	// kick off a worker
	// this would normally happen on a different machine ..
	ocrWorker, err := NewOcrRpcWorker(rabbitConfig)
	if err != nil {
		return err
	}
	ocrWorker.Run()
	return nil

}
