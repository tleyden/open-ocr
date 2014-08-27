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

	ocrRequest := OcrRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&ocrRequest)
	if err != nil {
		logg.LogError(err)
		http.Error(w, "Unable to unmarshal json", 500)
		return
	}

	ocrResult, err := HandleOcrRequest(ocrRequest, s.RabbitConfig)

	if err != nil {
		msg := "Unable to perform OCR decode.  Error: %v"
		errMsg := fmt.Sprintf(msg, err)
		logg.LogError(fmt.Errorf(errMsg))
		http.Error(w, errMsg, 500)
		return
	}

	logg.LogTo("OCR_HTTP", "ocrResult: %v", ocrResult)

	fmt.Fprintf(w, ocrResult.Text)

}

func HandleOcrRequest(ocrRequest OcrRequest, rabbitConfig RabbitConfig) (OcrResult, error) {

	switch ocrRequest.InplaceDecode {
	case true:
		// inplace decode: short circuit rabbitmq, and just call
		// ocr engine directly
		ocrEngine := NewOcrEngine(ocrRequest.EngineType)

		ocrResult, err := ocrEngine.ProcessRequest(ocrRequest)

		if err != nil {
			msg := "Error processing ocr request.  Error: %v"
			errMsg := fmt.Sprintf(msg, err)
			logg.LogError(fmt.Errorf(errMsg))
			return OcrResult{}, err
		}

		return ocrResult, nil
	default:
		// add a new job to rabbitmq and wait for worker to respond w/ result
		ocrClient, err := NewOcrRpcClient(rabbitConfig)
		if err != nil {
			logg.LogError(err)
			return OcrResult{}, err
		}

		ocrResult, err := ocrClient.DecodeImage(ocrRequest)

		if err != nil {
			logg.LogError(err)
			return OcrResult{}, err
		}

		return ocrResult, nil
	}

}
