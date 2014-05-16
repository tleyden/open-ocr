package ocrworker

import (
	"encoding/json"
	"fmt"
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

	var object interface{}

	if err := json.Unmarshal(b, &object); err != nil {
		return err
	}

	switch object.(type) {
	case int, int8, int16, int32, int64:
		if object, ok := object.(int); ok {
			*e = OcrEngineType(object)
		} else {
			return fmt.Errorf("Error unmarshaling OcrEngineType, not int")
		}

	case string:
		if object, ok := object.(string); ok {
			engineString := strings.ToUpper(object)
			switch engineString {
			case "TESSERACT":
				*e = ENGINE_TESSERACT
			case "MOCK":
				*e = ENGINE_MOCK
			default:
				logg.LogWarn("Unexpected OcrEngineType json: %v", engineString)
				*e = ENGINE_MOCK
			}

		} else {
			return fmt.Errorf("Error unmarshaling OcrEngineType, not string")
		}

	default:
		logg.LogWarn("Got unexpected type OcrEngineType json: %T", object)
		*e = ENGINE_MOCK
	}

	return
}
