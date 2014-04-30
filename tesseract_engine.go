package ocrworker

type TesseractEngine struct {
}

func (t TesseractEngine) ProcessImageUrl(imgUrl string) (OcrResult, error) {
	return OcrResult{}, nil
}
