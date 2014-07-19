package ocrworker

import (
	"errors"
	"fmt"
	"os"

	"github.com/GeertJohan/go.leptonica"
	"github.com/GeertJohan/go.tesseract"
	"github.com/couchbaselabs/logg"
)

const TESSERACT_MODEL_DIR = "/usr/local/share/tessdata"
const TESSERACT_LANG = "eng"

type GoTesseractEngine struct {
}

func (t GoTesseractEngine) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {

	ocrResult := OcrResult{Text: "Error"}
	err := errors.New("")

	if ocrRequest.ImgUrl != "" {
		ocrResult, err = t.ProcessImageUrl(ocrRequest.ImgUrl)
	} else {
		ocrResult, err = t.ProcessImageBytes(ocrRequest.ImgBytes)
	}

	return ocrResult, err

}

func (t GoTesseractEngine) ProcessImageBytes(imgBytes []byte) (OcrResult, error) {

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

func (t GoTesseractEngine) ProcessImageUrl(imgUrl string) (OcrResult, error) {

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

func (t GoTesseractEngine) processImageFile(tmpFileName string) (OcrResult, error) {

	tess, err := tesseract.NewTess(TESSERACT_MODEL_DIR, TESSERACT_LANG)
	if err != nil {
		return OcrResult{}, err
	}
	defer tess.Close()

	pix, err := leptonica.NewPixFromFile(tmpFileName)

	if err != nil {
		return OcrResult{}, err
	}
	defer pix.Close()

	// set the image to the tesseract instance
	tess.SetImagePix(pix)

	// retrieve text from the tesseract instance
	fmt.Println(tess.Text())

	return OcrResult{
		Text: tess.Text(),
	}, nil

}
