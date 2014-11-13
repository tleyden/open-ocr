package ocrworker

import (
	"encoding/json"
	"strings"

	"github.com/couchbaselabs/logg"
)

type OcrEngineType int

const (
	ENGINE_TESSERACT = OcrEngineType(iota)
	ENGINE_GO_TESSERACT
	ENGINE_MOCK
)

type OcrEngine interface {
	ProcessRequest(ocrRequest OcrRequest) (OcrResult, error)
}

func NewOcrEngine(engineType OcrEngineType) OcrEngine {
	switch engineType {
	case ENGINE_MOCK:
		return &MockEngine{}
	case ENGINE_TESSERACT:
		return &TesseractEngine{}
	}
	return nil
}

func (e OcrEngineType) String() string {
	switch e {
	case ENGINE_MOCK:
		return "ENGINE_MOCK"
	case ENGINE_TESSERACT:
		return "ENGINE_TESSERACT"
	case ENGINE_GO_TESSERACT:
		return "ENGINE_GO_TESSERACT"

	}
	return ""
}

func (e *OcrEngineType) UnmarshalJSON(b []byte) (err error) {

	var engineTypeStr string

	if err := json.Unmarshal(b, &engineTypeStr); err == nil {
		engineString := strings.ToUpper(engineTypeStr)
		switch engineString {
		case "TESSERACT":
			*e = ENGINE_TESSERACT
		case "GO_TESSERACT":
			*e = ENGINE_GO_TESSERACT
		case "MOCK":
			*e = ENGINE_MOCK
		default:
			logg.LogWarn("Unexpected OcrEngineType json: %v", engineString)
			*e = ENGINE_MOCK
		}
		return nil
	}

	// not a string .. maybe it's an int

	var engineTypeInt int
	if err := json.Unmarshal(b, &engineTypeInt); err == nil {
		*e = OcrEngineType(engineTypeInt)
		return nil
	} else {
		return err
	}

}
