package main

import (
	"fmt"

	"github.com/couchbaselabs/logg"
	"github.com/tleyden/open-ocr"
)

// This assumes that there is a rabbit mq running
// To test it, fire up a webserver and send it a curl request

func init() {
	logg.LogKeys["OCR"] = true
	logg.LogKeys["OCR_CLIENT"] = true
	logg.LogKeys["OCR_WORKER"] = true
	logg.LogKeys["PREPROCESSOR_WORKER"] = true
	logg.LogKeys["OCR_HTTP"] = true
	logg.LogKeys["OCR_TESSERACT"] = true
}

const CUSTOM_PREPROCESSOR_KEY = "custom-preprocessor"

type CustomPreprocessor struct {
}

func NewCustomPreprocessor() *CustomPreprocessor {
	return &CustomPreprocessor{}
}

func (c CustomPreprocessor) preprocess(ocrRequest *OcrRequest) error {
	// rotate the image
	// get bytes
	// ocrRequest.ImgBytes = <rotated bytes>
}

func main() {

	rabbitConfig := ocrworker.DefaultTestConfig()

	// inifinite loop, since sometimes worker <-> rabbitmq connection
	// gets broken.  see https://github.com/tleyden/open-ocr/issues/4
	for {
		logg.LogTo("PREPROCESSOR_WORKER", "Creating new Preprocessor Worker")

		ocrworker.RegisterPreprocessor(CUSTOM_PREPROCESSOR_KEY, NewCustomPreprocessor())

		preprocessorWorker, err := ocrworker.NewPreprocessorRpcWorker(
			rabbitConfig,
			CUSTOM_PREPROCESSOR_KEY,
		)
		if err != nil {
			logg.LogPanic("Could not create rpc worker: %v", err)
		}
		preprocessorWorker.Run()

		// this happens when connection is closed
		err = <-preprocessorWorker.Done
		logg.LogError(fmt.Errorf("Preprocessor Worker failed with error: %v", err))
	}

}
