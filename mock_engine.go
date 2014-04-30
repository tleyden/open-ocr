package ocrworker

type MockEngine struct {
}

func (m MockEngine) ProcessImageUrl(imgUrl string) (OcrResult, error) {
	return OcrResult{Text: "the lazy brown fox"}, nil
}
