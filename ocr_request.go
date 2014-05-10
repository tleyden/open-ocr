package ocrworker

type OcrRequest struct {
	ImgUrl     string        `json:"img_url"`
	EngineType OcrEngineType `json:"engine"`
}
