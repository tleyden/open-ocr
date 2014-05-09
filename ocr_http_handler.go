package ocrworker

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	"io/ioutil"
	"net/http"
)

type OcrHttpHandler struct {
}

func NewOcrHttpHandler(r RabbitConfig) *OcrHttpHandler {
	return &OcrHttpHandler{}
}

func (s *OcrHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	logg.LogTo("OCR_HTTP", "serveHttp called")
	defer req.Body.Close()
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logg.LogPanic("Unable to read request body")
	}
	bodyStr := string(bodyBytes)
	logg.LogTo("OCR_HTTP", "bodyStr: %v", bodyStr)
	fmt.Fprintf(w, "Hello, world")

}
