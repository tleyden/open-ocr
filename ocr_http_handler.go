package ocrworker

import (
	"fmt"
	"net/http"
)

type OcrHttpHandler struct {
}

func NewOcrHttpHandler(r RabbitConfig, e OcrEngineType) *OcrHttpHandler {
	return &OcrHttpHandler{}
}

func (s *OcrHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, world")
}
