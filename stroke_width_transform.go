package ocrworker

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/couchbaselabs/logg"
)

type StrokeWidthTransformer struct {
}

func (s StrokeWidthTransformer) preprocess(ocrRequest *OcrRequest) error {

	// write bytes to a temp file

	tmpFileNameInput, err := createTempFileName()
	tmpFileNameInput = fmt.Sprintf("%s.png", tmpFileNameInput)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileNameInput)

	tmpFileNameOutput, err := createTempFileName()
	tmpFileNameOutput = fmt.Sprintf("%s.png", tmpFileNameOutput)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileNameOutput)

	err = saveBytesToFileName(ocrRequest.ImgBytes, tmpFileNameInput)
	if err != nil {
		return err
	}

	// run DecodeText binary on it (if not in path, print warning and do nothing)
	darkOnLightSetting := s.extractDarkOnLightParam(*ocrRequest)
	logg.LogTo(
		"PREPROCESSOR_WORKER",
		"DetectText on %s -> %s with %s",
		tmpFileNameInput,
		tmpFileNameOutput,
		darkOnLightSetting,
	)
	out, err := exec.Command(
		"DetectText",
		tmpFileNameInput,
		tmpFileNameOutput,
		darkOnLightSetting,
	).CombinedOutput()
	if err != nil {
		logg.LogFatal("Error running command: %s.  out: %s", err, out)
	}
	logg.LogTo("PREPROCESSOR_WORKER", "output: %v", string(out))

	// read bytes from output file into ocrRequest.ImgBytes
	resultBytes, err := ioutil.ReadFile(tmpFileNameOutput)
	if err != nil {
		return err
	}

	ocrRequest.ImgBytes = resultBytes

	return nil

}

func (s StrokeWidthTransformer) extractDarkOnLightParam(ocrRequest OcrRequest) string {

	logg.LogTo("PREPROCESSOR_WORKER", "extract dark on light param")

	val := "1"

	preprocessorArgs := ocrRequest.PreprocessorArgs
	swtArgs := preprocessorArgs[PREPROCESSOR_STROKE_WIDTH_TRANSFORM]
	if swtArgs != nil {
		swtArg, ok := swtArgs.(string)
		if ok && (swtArg == "0" || swtArg == "1") {
			val = swtArg
		}
	}

	logg.LogTo("PREPROCESSOR_WORKER", "return val: %s", val)

	return val

}
