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
// TODO: update dockerfile to install go-tesseract package!
type TesseractEngineExec struct {
}

func (t TesseractEngineExec) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {
	return OcrResult{}, nil
}

func (t TesseractEngineExec) ProcessImageBytes(imgBytes []byte) (OcrResult, error) {

	tmpFileName, err := createTempFileName()
	if err != nil {
		return OcrResult{}, err
	}
	defer os.Remove(tmpFileName)

	// we have to write the contents of the image url to a temp
	// file, because the leptonica lib can't seem to handle byte arrays
	err = saveBytesToFileName(imgBytes, tmpFileName)
	if err != nil {
		return OcrResult{}, err
	}

	return t.processImageFile(tmpFileName)

}

func (t TesseractEngineExec) ProcessImageUrl(imgUrl string) (OcrResult, error) {

	logg.LogTo("OCR_TESSERACT", "ProcessImageUrl()")

	tmpFileName, err := createTempFileName()
	if err != nil {
		return OcrResult{}, err
	}
	defer os.Remove(tmpFileName)
	// we have to write the contents of the image url to a temp
	// file, because the leptonica lib can't seem to handle byte arrays
	err = saveUrlContentToFileName(imgUrl, tmpFileName)
	if err != nil {
		return OcrResult{}, err
	}

	return t.processImageFile(tmpFileName)

}

func (t TesseractEngineExec) processImageFile(inputFilename string) (OcrResult, error) {

	// give tesseract a unique output filename
	tmpOutFileBaseName, err := createTempFileName()
	if err != nil {
		logg.LogTo("OCR_TESSERACT", "Error creating tmp file: %v", err)
		return OcrResult{}, err
	}

	// the actual file it writes to will have a .txt extension
	tmpOutFileName := fmt.Sprintf("%s.txt", tmpOutFileBaseName)

	// delete output file when we are done
	defer os.Remove(tmpOutFileName)

	// exec tesseract
	cmd := exec.Command("tesseract", inputFilename, tmpOutFileBaseName)
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
