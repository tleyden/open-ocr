package ocrworker

type OcrEngineType int

const (
	ENGINE_TESSERACT = OcrEngineType(iota)
	ENGINE_MOCK
)
