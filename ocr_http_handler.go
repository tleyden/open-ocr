package ocrworker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbaselabs/logg"
)

type OcrHttpHandler struct {
	RabbitConfig RabbitConfig
}

func NewOcrHttpHandler(r RabbitConfig) *OcrHttpHandler {
	return &OcrHttpHandler{
		RabbitConfig: r,
	}
}

func (s *OcrHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	logg.LogTo("OCR_HTTP", "serveHttp called")
	defer req.Body.Close()

	ocrReq := OcrRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&ocrReq)
	if err != nil {
		logg.LogError(err)
		http.Error(w, "Unable to unmarshal json", 500)
		return
	}

	// TODO: call func HandleOcrRequest(ocrRequest OcrRequest, rabbitConfig RabbitConfig) (OcrResult, error) instead of
	// code below

	ocrClient, err := NewOcrRpcClient(s.RabbitConfig)
	if err != nil {
		logg.LogError(err)
		http.Error(w, "Unable to create rpc client", 500)
		return
	}

	decodeResult, err := ocrClient.DecodeImage(ocrReq)

	if err != nil {
		logg.LogError(err)
		http.Error(w, "Unable to perform OCR decode", 500)
		return
	}

	logg.LogTo("OCR_HTTP", "decodeResult: %v", decodeResult)

	logg.LogTo("OCR_HTTP", "ocrReq: %v", ocrReq)
	fmt.Fprintf(w, decodeResult.Text)

}
