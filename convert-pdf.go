package ocrworker

/*	Use this module if you want to call tesseract over
	pdfsandwich with an image as input file.
	Useful with big documents.

	Use cases:
	engine: tesseract with file_type: pdf and preprocessor: convert-pdf
	engine: sandwich  with file_type: [tif, png, jpg] and preprocessor: convert-pdf
*/

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	"io/ioutil"
	"os"
	"os/exec"
)

type ConvertPdf struct {
}

func (c ConvertPdf) preprocess(ocrRequest *OcrRequest) error {

	tmpFileNameInput, err := createTempFileName()
	tmpFileNameInput = fmt.Sprintf("%s.pdf", tmpFileNameInput)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileNameInput)

	tmpFileNameOutput, err := createTempFileName()
	tmpFileNameOutput = fmt.Sprintf("%s.tif", tmpFileNameOutput)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFileNameOutput)

	err = saveBytesToFileName(ocrRequest.ImgBytes, tmpFileNameInput)
	if err != nil {
		return err
	}

	logg.LogTo(
		"PREPROCESSOR_WORKER",
		"Convert PDF %s -> %s",
		tmpFileNameInput,
		tmpFileNameOutput,
	)
	/*
	   gs -dQUIET -o Antrag.tif -sDEVICE=tiffg4 Antrag.pdf
	*/
	var gsArgs []string
	gsArgs = append(gsArgs,
		"-dQUIET",
		"-dNOPAUSE",
		"-dBATCH",
		"-sOutputFile="+tmpFileNameOutput,
		"-sDEVICE=tiffg4",
		tmpFileNameInput,
	)
	logg.LogTo("PREPROCESSOR_WORKER", "output: %s", gsArgs)

	out, err := exec.Command("gs", gsArgs...).CombinedOutput()
	if err != nil {
		logg.LogFatal("Error running command: %s. out: %s", err, out)
	}
	logg.LogTo("PREPROCESSOR_WORKER", "output: %v", string(out))

	// read bytes from output file
	resultBytes, err := ioutil.ReadFile(tmpFileNameOutput)

	if err != nil {
		return err
	}
	ocrRequest.ImgBytes = resultBytes

	return nil
}
