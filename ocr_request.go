package ocrworker

import "fmt"

type OcrRequest struct {
	ImgUrl            string                 `json:"img_url"`
	EngineType        OcrEngineType          `json:"engine"`
	ImgBytes          []byte                 `json:"img_bytes"`
	PreprocessorChain []string               `json:"preprocessors"`
	PreprocessorArgs  map[string]interface{} `json:"preprocessor-args"`

	// TODO: need a way to have it generic here, but able to convert
	// into a TesseractEngineExec struct that is specific to TesserectExec

	EngineArgs interface{} `json:"engine_args"`
}

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

func (ocrRequest *OcrRequest) downloadImgUrl() error {

	bytes, err := url2bytes(ocrRequest.ImgUrl)
	if err != nil {
		return err
	}
	ocrRequest.ImgBytes = bytes
	ocrRequest.ImgUrl = ""
	return nil
}

func (o OcrRequest) String() string {
	return fmt.Sprintf("ImgUrl: %s, EngineType: %s, Preprocessors: %s", o.ImgUrl, o.EngineType, o.PreprocessorChain)
}
