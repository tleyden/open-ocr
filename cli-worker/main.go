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
	logg.LogKeys["OCR_HTTP"] = true
	logg.LogKeys["OCR_TESSERACT"] = true
}

func main() {

	noOpFlagFunc := ocrworker.NoOpFlagFunction()
	rabbitConfig := ocrworker.DefaultConfigFlagsOverride(noOpFlagFunc)

	// inifinite loop, since sometimes worker <-> rabbitmq connection
	// gets broken.  see https://github.com/tleyden/open-ocr/issues/4
	for {
		logg.LogTo("OCR_WORKER", "Creating new OCR Worker")
		ocrWorker, err := ocrworker.NewOcrRpcWorker(rabbitConfig)
		if err != nil {
			logg.LogPanic("Could not create rpc worker")
		}

		if err := ocrWorker.Run(); err != nil {
			logg.LogPanic("Error running worker: %v", err)
		}

		// this happens when connection is closed
		err = <-ocrWorker.Done
		logg.LogError(fmt.Errorf("OCR Worker failed with error: %v", err))
	}

}
