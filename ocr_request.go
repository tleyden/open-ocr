package ocrworker

type OcrRequest struct {
	ImgUrl            string        `json:"img_url"`
	EngineType        OcrEngineType `json:"engine"`
	ImgBytes          []byte        `json:"img_bytes"`
	PreprocessorChain []string      `json:"preprocessors"`
}

// 		ocrRequest.nextPreprocessor(c.rabbitConfig.RoutingKey),

// figure out the next pre-processor routing key to use (if any).
// if we have finished with the pre-processors, then use the processorRoutingKey
func (ocrRequest *OcrRequest) nextPreprocessor(processorRoutingKey string) string {
	if len(ocrRequest.PreprocessorChain) == 0 {
		return processorRoutingKey
	} else {
		var x string
		s := ocrRequest.PreprocessorChain
		x, s = s[len(s)-1], s[:len(s)-1]
		ocrRequest.PreprocessorChain = s
		return x
	}

}
