package ocrworker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/couchbaselabs/logg"
)

type SwtOcrRequest struct {
	OcrRequest
	PreprocessorArgs []string `json:"preprocessor-args"`
}

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

	defaultVal := "1" // dark text on light background

	ocrRequestJson, err := json.Marshal(ocrRequest)
	if err != nil {
		logg.LogTo("PREPROCESSOR_WORKER", "got error: %v", err)
		logg.LogError(err)
		return defaultVal
	}

	swtOcrRequest := SwtOcrRequest{}
	err = json.Unmarshal(ocrRequestJson, &swtOcrRequest)
	if err != nil {
		logg.LogTo("PREPROCESSOR_WORKER", "got error: %v", err)
		logg.LogError(err)
		return defaultVal
	}

	if len(swtOcrRequest.PreprocessorArgs) > 0 {
		val := swtOcrRequest.PreprocessorArgs[0]
		logg.LogTo("PREPROCESSOR_WORKER", "dark on light param: %q", val)
		return val
	}

	logg.LogTo("PREPROCESSOR_WORKER", "return default val")

	return defaultVal

}
