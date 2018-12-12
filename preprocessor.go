package ocrworker

const PreprocessorIdentity = "identity"
const PreprocessorStrokeWidthTransform = "stroke-width-transform"
const PreprocessorConvertPdf = "convert-pdf"

type Preprocessor interface {
	preprocess(ocrRequest *OcrRequest) error
}

type IdentityPreprocessor struct {
}

func (i IdentityPreprocessor) preprocess(ocrRequest *OcrRequest) error {
	return nil
}
