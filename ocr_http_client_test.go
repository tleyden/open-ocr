package ocrworker

import (
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"github.com/tleyden/fakehttp"
	"testing"
)

func TestDecodeImageUrl(t *testing.T) {

	port := 8083
	fakeDecodedOcr := "fake ocr"
	sourceServer := fakehttp.NewHTTPServerWithPort(port)
	sourceServer.Start()
	headers := map[string]string{"Content-Type": "text/plain"}
	sourceServer.Response(200, headers, fakeDecodedOcr)

	openOcrUrl := fmt.Sprintf("http://localhost:%d", port)
	openOcrClient := NewHttpClient(openOcrUrl)
	attachmentUrl := "http://fake.io/a.png"
	ocrDecoded, err := openOcrClient.DecodeImageUrl(attachmentUrl, ENGINE_TESSERACT)
	assert.True(t, err == nil)
	assert.Equals(t, ocrDecoded, fakeDecodedOcr)

}
