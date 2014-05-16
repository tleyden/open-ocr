package ocrworker

import (
	"encoding/json"
	"github.com/couchbaselabs/logg"
	"strings"
)

type OcrEngineType int

const (
	ENGINE_TESSERACT = OcrEngineType(iota)
	ENGINE_MOCK
)

type OcrEngine interface {
	ProcessImageUrl(imgUrl string) (OcrResult, error)
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

func (e *OcrEngineType) UnmarshalJSON(b []byte) (err error) {

	var engineTypeStr string

	if err := json.Unmarshal(b, &engineTypeStr); err == nil {
		logg.LogTo("OCR", "its a string")
		engineString := strings.ToUpper(engineTypeStr)
		switch engineString {
		case "TESSERACT":
			*e = ENGINE_TESSERACT
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
