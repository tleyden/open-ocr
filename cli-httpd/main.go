package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/couchbaselabs/logg"
	"github.com/tleyden/open-ocr"
)

// This assumes that there is a worker running
// To test it:
// curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://localhost:8081/img","engine":0}' http://localhost:8081/ocr

func init() {
	logg.LogKeys["OCR"] = true
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

	// any requests to root, just redirect to main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		text := `<h1>OpenOCR is running!<h1> Need <a href="http://www.openocr.net">docs</a>?`
		fmt.Fprintf(w, text)
	})

	http.Handle("/ocr", ocrworker.NewOcrHttpHandler(rabbitConfig))

	// add a handler which will take both the JSON and the image data
	// in the same multipart/related request
	http.HandleFunc("/ocr-file-upload", func(w http.ResponseWriter, r *http.Request) {

		logg.LogTo("OCR_HTTP", "request to ocr-file-upload")

		switch r.Method {
		case "POST":
			h := r.Header.Get("Content-Type")
			logg.LogTo("OCR_HTTP", "content type: %v", h)

			contentType, attrs, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			logg.LogTo("OCR_HTTP", "content type: %v", contentType)

			if !strings.HasPrefix(h, "multipart/related") {
				http.Error(w, "Expected multipart related", 500)
				return
			}

			reader := multipart.NewReader(r.Body, attrs["boundary"])

			ocrReq := OcrRequest{}

			for {
				part, err := reader.NextPart()
				defer part.Close()

				if err == io.EOF {
					break
				}
				var body Body
				contentTypeOuter := mainPart.Header["Content-Type"][0]
				contentType, attrs, _ := mime.ParseMediaType(contentTypeOuter)
				switch contentType {
				case "application/json":
					decoder := json.NewDecoder(part)
					err := decoder.Decode(&ocrReq)
					if err != nil {
						logg.LogError(err)
						http.Error(w, "Unable to unmarshal json", 500)
						return
					}

				default:
					if !strings.HasPrefix(contentType, "image") {

						http.Error(w, "Expected content-type to start with image/", 500)
						return
					}

					// dump part to output (for now ..)

					partContents, err := ioutil.ReadAll(part)
					if err != nil {
						logg.LogTo("OCR_HTTP", "failed to read mime part: %v", part)
						return err
					}
					logg.LogTo("OCR_HTTP", "partContents: %v", partContents)

				}

			}

		}
	})

	// add a handler to serve up an image from the filesystem.
	// ignore this, was just something for testing ..
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../refactoring.png")
	})

	listenAddr := fmt.Sprintf(":%d", http_port)

	logg.LogTo("OCR_HTTP", "Starting listener on %v", listenAddr)
	logg.LogError(http.ListenAndServe(listenAddr, nil))

}
