package ocrworker

import (
	"encoding/json"
	"testing"

	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
)

func TestOcrEngineTypeJson(t *testing.T) {

	testJson := `{"img_url":"foo", "engine":"tesseract"}`
	ocrRequest := OcrRequest{}
	err := json.Unmarshal([]byte(testJson), &ocrRequest)
	if err != nil {
		logg.LogError(err)
	}
	assert.True(t, err == nil)
	assert.Equals(t, ocrRequest.EngineType, ENGINE_TESSERACT)
	logg.LogTo("TEST", "ocrRequest: %v", ocrRequest)

}
