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

	rabbitConfig := ocrworker.DefaultTestConfig()

	ocrWorker, err := ocrworker.NewOcrRpcWorker(rabbitConfig)
	if err != nil {
		logg.LogPanic("Could not create rpc worker")
	}
	ocrWorker.Run()

	select {}

}
