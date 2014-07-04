package ocrworker

import (
	"testing"

	"io/ioutil"

	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
)

func TestTesseractEngineExecWithRequest(t *testing.T) {

	engine := TesseractEngineExec{}
	bytes, err := ioutil.ReadFile("docs/testimage.png")

	ocrRequest := OcrRequest{
		ImgBytes:   bytes,
		EngineType: ENGINE_TESSERACT_EXEC,
	}

	assert.True(t, err == nil)
	result, err := engine.ProcessRequest(ocrRequest)
	assert.True(t, err == nil)
	logg.LogTo("TEST", "result: %v", result)

}

func TestTesseractEngineExecWithFile(t *testing.T) {

	engine := TesseractEngineExec{}
	result, err := engine.processImageFile("docs/testimage.png")
	assert.True(t, err == nil)
	logg.LogTo("TEST", "result: %v", result)

}

func TestTesseractEngineExecWithBytes(t *testing.T) {

	engine := TesseractEngineExec{}
	bytes, err := ioutil.ReadFile("docs/testimage.png")
	assert.True(t, err == nil)
	result, err := engine.ProcessImageBytes(bytes)
	assert.True(t, err == nil)
	logg.LogTo("TEST", "result: %v", result)

}
