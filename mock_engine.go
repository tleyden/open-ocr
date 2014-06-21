package ocrworker

const MOCK_ENGINE_RESPONSE = "mock engine decoder response"

type MockEngine struct {
}

func (m MockEngine) ProcessImageUrl(imgUrl string) (OcrResult, error) {
	return OcrResult{Text: MOCK_ENGINE_RESPONSE}, nil
}

func (m MockEngine) ProcessImageBytes(imgBytes []byte) (OcrResult, error) {
	return OcrResult{Text: MOCK_ENGINE_RESPONSE}, nil
}
