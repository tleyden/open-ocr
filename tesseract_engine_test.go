package ocrworker

import (
	"testing"

	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
)

func TestTesseractEngine(t *testing.T) {

	engine := TesseractEngine{}
	result, err := engine.processImageFile("docs/testimage.png")
	assert.True(t, err == nil)
	logg.LogTo("TEST", "result: %v", result)

}
