package ocrworker

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
}
