package ocrworker

import (
	"fmt"
	"github.com/GeertJohan/go.leptonica"
	"github.com/GeertJohan/go.tesseract"
	"github.com/couchbaselabs/logg"
	"github.com/nu7hatch/gouuid"
	"os"
	"path/filepath"
)

const TESSERACT_MODEL_DIR = "/usr/local/share/tessdata"
const TESSERACT_LANG = "eng"

type TesseractEngine struct {
}

func (t TesseractEngine) ProcessImageUrl(imgUrl string) (OcrResult, error) {

	// TODO: plug in tesseraact

	// create new tess instance and point it to the tessdata location. Set language to english.

	tess, err := tesseract.NewTess(TESSERACT_MODEL_DIR, TESSERACT_LANG)
	if err != nil {
		return OcrResult{}, err
	}
	defer tess.Close()

	tmpFileName, err := createTempFileName()
	if err != nil {
		return OcrResult{}, err
	}

	// we have to write the contents of the image url to a temp
	// file, because the leptonica lib can't seem to handle byte arrays
	err = saveUrlContentToFileName(imgUrl, tmpFileName)
	if err != nil {
		return OcrResult{}, err
	}

	logg.LogTo("OCR_TESSERACT", "saved %v to %v", imgUrl, tmpFileName)

	pix, err := leptonica.NewPixFromFile(tmpFileName)
	if err != nil {
		return OcrResult{}, err
	}

	// set the image to the tesseract instance
	tess.SetImagePix(pix)

	// retrieve text from the tesseract instance
	fmt.Println(tess.Text())

	return OcrResult{
		Text: tess.Text(),
	}, nil
}

func createTempFileName() (string, error) {
	tempDir := os.TempDir()
	uuidRaw, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	uuidStr := uuidRaw.String()
	return filepath.Join(tempDir, uuidStr), nil
}
