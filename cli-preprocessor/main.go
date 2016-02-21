package main

import (
	"flag"
	"fmt"

	"github.com/couchbaselabs/logg"
	ocrworker "github.com/tleyden/open-ocr"
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

func main() {

	var preprocessor string
	flagFunc := func() {
		flag.StringVar(
			&preprocessor,
			"preprocessor",
			"identity",
			"The preprocessor to use, eg, stroke-width-transform",
		)

	}

	rabbitConfig := ocrworker.DefaultConfigFlagsOverride(flagFunc)

	// inifinite loop, since sometimes worker <-> rabbitmq connection
	// gets broken.  see https://github.com/tleyden/open-ocr/issues/4
	for {
		logg.LogTo("PREPROCESSOR_WORKER", "Creating new Preprocessor Worker")
		preprocessorWorker, err := ocrworker.NewPreprocessorRpcWorker(
			rabbitConfig,
			preprocessor,
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
