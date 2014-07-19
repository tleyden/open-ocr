package ocrworker

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/couchbaselabs/logg"
)

// This variant of the TesseractEngine calls tesseract via exec rather
// than go.tesseract.
// TODO: update dockerfile to install tesseract command line package!
type TesseractEngineExec struct {
}

type TesseractEngineExecArgs struct {
	configVars  map[string]string `json:"config_vars"`
	pageSegMode string            `json:"psm"`
}

func NewTesseractEngineExecArgs(ocrRequest OcrRequest) (*TesseractEngineExecArgs, error) {

	engineArgs := &TesseractEngineExecArgs{}

	if ocrRequest.EngineArgs == nil {
		return engineArgs, nil
	}

	// config vars
	configVarsMapInterfaceOrig := ocrRequest.EngineArgs["config_vars"]

	if configVarsMapInterfaceOrig != nil {

		logg.LogTo("OCR_TESSERACT", "got configVarsMap: %v type: %T", configVarsMapInterfaceOrig, configVarsMapInterfaceOrig)

		configVarsMapInterface := configVarsMapInterfaceOrig.(map[string]interface{})

		configVarsMap := make(map[string]string)
		for k, v := range configVarsMapInterface {
			v, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("Could not convert configVar into string: %v", v)
			}
			configVarsMap[k] = v
		}

		engineArgs.configVars = configVarsMap

	}

	// page seg mode
	pageSegMode := ocrRequest.EngineArgs["psm"]
	if pageSegMode != nil {
		pageSegModeStr, ok := pageSegMode.(string)
		if !ok {
			return nil, fmt.Errorf("Could not convert psm into string: %v", pageSegMode)
		}
		engineArgs.pageSegMode = pageSegModeStr
	}
	return engineArgs, nil

}

// return a slice that can be passed to tesseract binary as command line
// args, eg, ["-c", "tessedit_char_whitelist=0123456789", "-c", "foo=bar"]
func (t TesseractEngineExecArgs) Export() []string {
	result := []string{}
	for k, v := range t.configVars {
		result = append(result, "-c")
		keyValArg := fmt.Sprintf("%s=%s", k, v)
		result = append(result, keyValArg)
	}
	if t.pageSegMode != "" {
		result = append(result, "-psm")
		result = append(result, t.pageSegMode)
	}
	return result
}

func (t TesseractEngineExec) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {

	tmpFileName, err := func() (string, error) {
		if ocrRequest.ImgUrl != "" {
			return t.tmpFileFromImageUrl(ocrRequest.ImgUrl)
		} else {
			return t.tmpFileFromImageBytes(ocrRequest.ImgBytes)
		}

	}()

	if err != nil {
		logg.LogTo("OCR_TESSERACT", "error getting tmpFileName")
		return OcrResult{}, err
	}

	defer os.Remove(tmpFileName)

	engineArgs, err := NewTesseractEngineExecArgs(ocrRequest)
	if err != nil {
		logg.LogTo("OCR_TESSERACT", "error getting engineArgs")
		return OcrResult{}, err
	}

	ocrResult, err := t.processImageFile(tmpFileName, *engineArgs)

	return ocrResult, err

}

func (t TesseractEngineExec) tmpFileFromImageBytes(imgBytes []byte) (string, error) {

	tmpFileName, err := createTempFileName()
	if err != nil {
		return "", err
	}

	// we have to write the contents of the image url to a temp
	// file, because the leptonica lib can't seem to handle byte arrays
	err = saveBytesToFileName(imgBytes, tmpFileName)
	if err != nil {
		return "", err
	}

	return tmpFileName, nil

}

func (t TesseractEngineExec) tmpFileFromImageUrl(imgUrl string) (string, error) {

	tmpFileName, err := createTempFileName()
	if err != nil {
		return "", err
	}
	// we have to write the contents of the image url to a temp
	// file, because the leptonica lib can't seem to handle byte arrays
	err = saveUrlContentToFileName(imgUrl, tmpFileName)
	if err != nil {
		return "", err
	}

	return tmpFileName, nil

}

func (t TesseractEngineExec) processImageFile(inputFilename string, engineArgs TesseractEngineExecArgs) (OcrResult, error) {

	// if the input filename is /tmp/ocrimage, set the output file basename
	// to /tmp/ocrimage as well, which will produce /tmp/ocrimage.txt output
	tmpOutFileBaseName := inputFilename

	// the actual file it writes to will have a .txt extension
	tmpOutFileName := fmt.Sprintf("%s.txt", inputFilename)

	// delete output file when we are done
	defer os.Remove(tmpOutFileName)

	// build args array
	cflags := engineArgs.Export()
	cmdArgs := []string{inputFilename, tmpOutFileBaseName}
	cmdArgs = append(cmdArgs, cflags...)
	logg.LogTo("OCR_TESSERACT", "cmdArgs: %v", cmdArgs)

	// exec tesseract
	cmd := exec.Command("tesseract", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logg.LogTo("OCR_TESSERACT", "Error exec tesseract: %v %v", err, string(output))
		return OcrResult{}, err
	}

	// get data from outfile
	outBytes, err := ioutil.ReadFile(tmpOutFileName)
	if err != nil {
		logg.LogTo("OCR_TESSERACT", "Error getting data from out file: %v", err)
		return OcrResult{}, err
	}

	return OcrResult{
		Text: string(outBytes),
	}, nil

}
