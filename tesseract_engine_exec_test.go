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

	cFlags := make(map[string]string)
	cFlags["tessedit_char_whitelist"] = "0123456789"
	engineArgs := TesseractEngineExecArgs{
		cFlags: cFlags,
	}

	ocrRequest := OcrRequest{
		ImgBytes:   bytes,
		EngineType: ENGINE_TESSERACT_EXEC,
		EngineArgs: engineArgs,
	}

	assert.True(t, err == nil)
	result, err := engine.ProcessRequest(ocrRequest)
	assert.True(t, err == nil)
	logg.LogTo("TEST", "result: %v", result)

}

func TestTesseractEngineExecWithFile(t *testing.T) {

	engine := TesseractEngineExec{}
	engineArgs := TesseractEngineExecArgs{}
	result, err := engine.processImageFile("docs/testimage.png", engineArgs)
	assert.True(t, err == nil)
	logg.LogTo("TEST", "result: %v", result)

}
