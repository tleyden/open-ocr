package ocrworker

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io"
	// "reflect"
	// "bytes"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	// "net/http/httputil"
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
	// logg.LogTo("OCR_HTTP", "headers:%v:", req.Header)
	// logg.LogTo("OCR_HTTP", "body type:%v:", reflect.TypeOf(req.Body))
	// logg.LogTo("OCR_HTTP", "body cl:%d:", req.ContentLength)
	// body, _ := ioutil.ReadAll(req.Body)
	// logg.LogTo("OCR_HTTP", "Body:%s:", body)

	ocrReq := OcrRequest{}

	switch req.Method {
	case "POST":
		h := req.Header.Get("Content-Type")
		logg.LogTo("OCR_HTTP", "content type: %v", h)

		contentType, attrs, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
		logg.LogTo("OCR_HTTP", "content type:%v:attrs:%v:", contentType, attrs)

		if !strings.HasPrefix(h, "multipart/related") {
			return ocrReq, fmt.Errorf("Expected multipart related")
		}

		reader := multipart.NewReader(req.Body, attrs["boundary"])
		logg.LogTo("OCR_HTTP", "got a reader:boundary:%s:", attrs["boundary"])

		for {

			logg.LogTo("OCR_HTTP", "in loop")
			part, err := reader.NextPart()
			logg.LogTo("OCR_HTTP", "got part:%s:err:%v:", part, err)

			if err == io.EOF {
				logg.LogTo("OCR_HTTP", "break out of loop")
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
				logg.LogTo("OCR_HTTP", "decoded req args:EngineArgs:%v:PreprocessorArgs:%v:PreprocessorChain:%v:", ocrReq.EngineArgs, ocrReq.PreprocessorArgs, ocrReq.PreprocessorChain)
				part.Close()
			default:
				if !strings.HasPrefix(contentType, "image") {
					return ocrReq, fmt.Errorf("Expected content-type: image/*")
				}

				partContents, err := ioutil.ReadAll(part)
				if err != nil {
					return ocrReq, fmt.Errorf("Failed to read mime part: %v", err)
				}

				// fmt.Printf("encoded image:%q:\n", partContents)
				buf := make([]byte, req.ContentLength)
				bytesRead, err := base64.StdEncoding.Decode(buf, partContents)
				// if err != nil {
				// 	fmt.Println("error:", err)
				// 	return ocrReq, fmt.Errorf("Unable to decode image: %s", err)
				// }
				if bytesRead == 0 && err != nil {
			        // log.Fatal(err)
					fmt.Println("error:", err)
					return ocrReq, fmt.Errorf("Unable to decode image: %s", err)
				}
				fmt.Printf("decoded:bytesRead:%d:\n", bytesRead)

				ocrReq.ImgBytes = buf[:bytesRead]
				logg.LogTo("OCR_HTTP", "final req args:EngineArgs:%v:PreprocessorArgs:%v:PreprocessorChain:%v:", ocrReq.EngineArgs, ocrReq.PreprocessorArgs, ocrReq.PreprocessorChain)
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

	// its gone when it gets here
	logg.LogTo("OCR_HTTP", "sending req args:EngineArgs:%v:PreprocessorArgs:%v:PreprocessorChain:%v:", ocrRequest.EngineArgs, ocrRequest.PreprocessorArgs, ocrRequest.PreprocessorChain)
	// logg.LogTo("OCR_HTTP", "sending ocrRequest: %v", ocrRequest)

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
