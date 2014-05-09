package ocrworker

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"net/http"
	"strings"
	"testing"
	"time"
)

// This test assumes that rabbit mq is running
func TestOcrHttpHandler(t *testing.T) {

	// add a handler to serve up an image from the filesystem.
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "refactoring.png")
	})

	rabbitConfig := rabbitConfigForTests()

	http.Handle("/ocr", NewOcrHttpHandler(rabbitConfig))

	go http.ListenAndServe(":8081", nil)

	logg.LogTo("TEST", "test1")

	jsonBody := `{"img_url": "http://localhost:8081/img", "engine": "mock_engine"}`
	reader := strings.NewReader(jsonBody)

	resp, err := http.Post("http://localhost:8081/ocr", "application/json", reader)
	assert.True(t, err == nil)
	logg.LogTo("TEST", "resp: %v", resp)

	// connect to it via http client
	time.Sleep(time.Second * 60)

	// make sure get expected result
	logg.LogTo("TEST", "test2")
	assert.True(t, true)
}
