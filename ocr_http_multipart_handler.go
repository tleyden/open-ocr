package ocrworker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/couchbaselabs/logg"
)

type OcrHttpMultipartHandler struct {
	RabbitConfig RabbitConfig
}

func NewOcrHttpMultipartHandler(r RabbitConfig) *OcrHttpMultipartHandler {
	return &OcrHttpMultipartHandler{
		RabbitConfig: r,
	}
}

func (s *OcrHttpMultipartHandler) extractParts(req *http.Request) (OcrRequest, error) {

	logg.LogTo("OCR_HTTP", "request to ocr-file-upload")
	ocrReq := OcrRequest{}

	switch req.Method {
	case "POST":
		h := req.Header.Get("Content-Type")
		logg.LogTo("OCR_HTTP", "content type: %v", h)

		contentType, attrs, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
		logg.LogTo("OCR_HTTP", "content type: %v", contentType)

		if !strings.HasPrefix(h, "multipart/related") {
			return ocrReq, fmt.Errorf("Expected multipart related")
		}

		reader := multipart.NewReader(req.Body, attrs["boundary"])

		for {

			part, err := reader.NextPart()

			if err == io.EOF {
				break
			}
			contentTypeOuter := part.Header["Content-Type"][0]
			contentType, attrs, _ := mime.ParseMediaType(contentTypeOuter)

			logg.LogTo("OCR_HTTP", "attrs: %v", attrs)

			switch contentType {
			case "application/json":
				decoder := json.NewDecoder(part)
				err := decoder.Decode(&ocrReq)
				if err != nil {
					return ocrReq, fmt.Errorf("Unable to unmarshal json: %s", err)
				}
				part.Close()
			default:
				if !strings.HasPrefix(contentType, "image") {
					return ocrReq, fmt.Errorf("Expected content-type: image/*")
				}

				partContents, err := ioutil.ReadAll(part)
				if err != nil {
					return ocrReq, fmt.Errorf("Failed to read mime part: %v", err)
				}

				ocrReq.ImgBytes = partContents
				return ocrReq, nil

			}

		}

		return ocrReq, fmt.Errorf("Didn't expect to get this far")

	default:
		return ocrReq, fmt.Errorf("This endpoint only accepts POST requests")
	}

}

func (s *OcrHttpMultipartHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	ocrRequest, err := s.extractParts(req)
	if err != nil {
		logg.LogError(err)
		errStr := fmt.Sprintf("Error extracting multipart/related parts: %v", err)
		http.Error(w, errStr, 500)
		return
	}

	logg.LogTo("OCR_HTTP", "ocrRequest: %v", ocrRequest)

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
