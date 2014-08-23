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

func (s *OcrHttpMultipartHandler) extractParts(req *http.Request) (OcrRequest, *multipart.Part, error) {

	logg.LogTo("OCR_HTTP", "request to ocr-file-upload")
	ocrReq := OcrRequest{}
	// var imagePart *multipart.Part

	switch req.Method {
	case "POST":
		h := req.Header.Get("Content-Type")
		logg.LogTo("OCR_HTTP", "content type: %v", h)

		contentType, attrs, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
		logg.LogTo("OCR_HTTP", "content type: %v", contentType)

		if !strings.HasPrefix(h, "multipart/related") {
			return ocrReq, nil, fmt.Errorf("Expected multipart related")
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
					return ocrReq, nil, fmt.Errorf("Unable to unmarshal json: %s", err)
				}
				part.Close()
			default:
				if !strings.HasPrefix(contentType, "image") {

					return ocrReq, nil, fmt.Errorf("Expected content-type: image/*")
				}

				// hack: this forces it to come in the order:
				// json / image, which I was trying to avoid.
				// was getting EOF when saving part and returning
				// TODO: read into []byte and return that
				return ocrReq, part, nil

			}

		}

		return ocrReq, nil, fmt.Errorf("Didn't expect to get this far")

	default:
		return ocrReq, nil, fmt.Errorf("This endpoint only accepts POST requests")
	}

}

func (s *OcrHttpMultipartHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	ocrRequest, imagePartReader, err := s.extractParts(req)
	if err != nil {
		logg.LogError(err)
		errStr := fmt.Sprint("%v", err)
		http.Error(w, errStr, 500)
		return
	}

	logg.LogTo("OCR_HTTP", "ocrRequest: %v", ocrRequest)

	// dump part to output (for now ..)

	partContents, err := ioutil.ReadAll(imagePartReader)
	if err != nil {
		logg.LogTo("OCR_HTTP", "Failed to read mime part: %v", err)
		http.Error(w, "Failed to read mime part", 500)
		return
	}
	ocrRequest.ImgBytes = partContents

	ocrClient, err := NewOcrRpcClient(s.RabbitConfig)
	if err != nil {
		logg.LogError(err)
		http.Error(w, "Unable to create rpc client", 500)
		return
	}

	decodeResult, err := ocrClient.DecodeImage(ocrRequest)

	if err != nil {
		logg.LogError(err)
		http.Error(w, "Unable to perform OCR decode", 500)
		return
	}

	logg.LogTo("OCR_HTTP", "decodeResult: %v", decodeResult)

	logg.LogTo("OCR_HTTP", "ocrReq: %v", ocrRequest)
	fmt.Fprintf(w, decodeResult.Text)

	/*

		                OLD CODE -- some of this still needs to be moved up

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
	*/

}
