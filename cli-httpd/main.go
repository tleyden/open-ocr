package main

import (
	"flag"
	"fmt"
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/open-ocr"
	"net/http"
)

// This assumes that there is a worker running
// To test it:
// curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://localhost:8081/img","engine":0}' http://localhost:8081/ocr

func init() {
	logg.LogKeys["OCR_CLIENT"] = true
	logg.LogKeys["OCR_WORKER"] = true
	logg.LogKeys["OCR_HTTP"] = true
	logg.LogKeys["OCR_TESSERACT"] = true
}

func main() {

	var http_port int
	flagFunc := func() {
		flag.IntVar(
			&http_port,
			"http_port",
			8080,
			"The http port to listen on, eg, 8081",
		)

	}
	rabbitConfig := ocrworker.DefaultConfigFlagsOverride(flagFunc)

	// add a handler to serve up an image from the filesystem.
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../refactoring.png")
	})

	http.Handle("/ocr", ocrworker.NewOcrHttpHandler(rabbitConfig))

	listenAddr := fmt.Sprintf(":%d", http_port)

	logg.LogTo("OCR_HTTP", "Starting listener on %v", listenAddr)
	logg.LogError(http.ListenAndServe(listenAddr, nil))

}
