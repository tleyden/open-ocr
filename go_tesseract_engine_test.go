package ocrworker

import (
	"testing"

	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
)

func TestGoTesseractEngine(t *testing.T) {

	engine := GoTesseractEngine{}
	result, err := engine.processImageFile("docs/testimage.png")
	assert.True(t, err == nil)
	logg.LogTo("TEST", "result: %v", result)

}
