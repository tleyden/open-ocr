package ocrworker

import (
	"encoding/json"
	"testing"

	"github.com/couchbaselabs/go.assert"
)

func TestParamExtraction(t *testing.T) {

	testJson := `{"img_url":"foo", "engine":"tesseract", "preprocessor-args":{"stroke-width-transform":"0"}}`
	ocrRequest := OcrRequest{}
	err := json.Unmarshal([]byte(testJson), &ocrRequest)
	assert.True(t, err == nil)

	swt := StrokeWidthTransformer{}
	param := swt.extractDarkOnLightParam(ocrRequest)
	assert.Equals(t, param, "0")

}

func TestParamExtractionNegative(t *testing.T) {

	testJson := `{"img_url":"foo", "engine":"tesseract"}`
	ocrRequest := OcrRequest{}
	err := json.Unmarshal([]byte(testJson), &ocrRequest)
	assert.True(t, err == nil)

	swt := StrokeWidthTransformer{}
	param := swt.extractDarkOnLightParam(ocrRequest)
	assert.Equals(t, param, "1")

}
