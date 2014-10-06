package ocrworker

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"os"
	"os/exec"

	"github.com/couchbaselabs/logg"
)

// This variant of the TesseractEngine calls tesseract via exec
type TesseractEngine struct {
}

type TesseractEngineArgs struct {
	configVars  map[string]string `json:"config_vars"`
	pageSegMode string            `json:"psm"`
	lang        string            `json:"lang"`
}

func NewTesseractEngineArgs(ocrRequest OcrRequest) (*TesseractEngineArgs, error) {

	engineArgs := &TesseractEngineArgs{}

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

	// language
	lang := ocrRequest.EngineArgs["lang"]
	if lang != nil {
		langStr, ok := lang.(string)
		if !ok {
			return nil, fmt.Errorf("Could not convert lang into string: %v", lang)
		}
		engineArgs.lang = langStr
	}

	return engineArgs, nil

}

// return a slice that can be passed to tesseract binary as command line
// args, eg, ["-c", "tessedit_char_whitelist=0123456789", "-c", "foo=bar"]
func (t TesseractEngineArgs) Export() []string {
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
	if t.lang != "" {
		result = append(result, "-l")
		result = append(result, t.lang)
	}

	return result
}

func (t TesseractEngine) ProcessRequest(ocrRequest OcrRequest) (OcrResult, error) {

	logg.LogTo("OCR_TESSERACT", "doing ocr:EngineArgs:%v:PreprocessorArgs:%v:PreprocessorChain:%v:", ocrRequest.EngineArgs, ocrRequest.PreprocessorArgs, ocrRequest.PreprocessorChain)

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

	engineArgs, err := NewTesseractEngineArgs(ocrRequest)
	if err != nil {
		logg.LogTo("OCR_TESSERACT", "error getting engineArgs")
		return OcrResult{}, err
	}

	logg.LogTo("OCR_TESSERACT", "doing ocr:engineArgs:%v:", engineArgs)
	ocrResult, err := t.processImageFile(tmpFileName, *engineArgs)

	return ocrResult, err

}

func (t TesseractEngine) tmpFileFromImageBytes(imgBytes []byte) (string, error) {

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

func (t TesseractEngine) tmpFileFromImageUrl(imgUrl string) (string, error) {

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

func (t TesseractEngine) processImageFile(inputFilename string, engineArgs TesseractEngineArgs) (OcrResult, error) {

	// if the input filename is /tmp/ocrimage, set the output file basename
	// to /tmp/ocrimage as well, which will produce /tmp/ocrimage.txt output
	tmpOutFileBaseName := inputFilename

	// the actual file it writes to will have a .txt extension
	tmpOutFileName := fmt.Sprintf("%s.txt", inputFilename)

	// delete output file when we are done
	if engineArgs.configVars["tessedit_write_images"] == "T" {

		// engineArgs.configVars
		// logg.LogTo("OCR_TESSERACT", "configVars: %v", engineArgs.configVars["tessedit_write_images"])

	} else {

		defer os.Remove(tmpOutFileName)

	}

	// /var/folders/9q/3j0vn7j52qd78qv49ldsw4cj6wf43f/T/445f10ac-1c05-49a9-4845-e06acf060ef3 
	// /var/folders/9q/3j0vn7j52qd78qv49ldsw4cj6wf43f/T/445f10ac-1c05-49a9-4845-e06acf060ef3 
	// -c language_model_penalty_non_dict_word=1 
	// -c language_model_penalty_non_freq_dict_word=1 
	// -c load_system_dawg=0 
	// -c tessedit_char_whitelist=ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789- 
	// -c tessedit_write_images=T 
	// -psm 1

	// build args array
	cflags := engineArgs.Export()

	//I guess tesseract adds a .txt to the inputFilename for the results
	cmdArgs := []string{inputFilename, tmpOutFileBaseName}
	cmdArgs = append(cmdArgs, cflags...)
	// cmdArgs2 := make(map[string]string)

	//  TESSDATA_PREFIX=/Users/bmcquee/git/openalpr/runtime_data/ocr/ alpr --config ~/git/openalpr/config/openalpr.conf  saab_ny.jpg -j
	logg.LogTo("OCR_TESSERACT", "cmdArgs: %v", cmdArgs)

	// exec tesseract
	// cmd := exec.Command("tesseract", cmdArgs...)
	// cmd := exec.Command("alpr", "--config", "/Users/bmcquee/git/openalpr/config/openalpr.conf", "-j", "/Users/bmcquee/Downloads/tonfiniti.jpg")
	cmd := exec.Command("alpr", "--config", "/Users/bmcquee/git/openalpr/config/openalpr.conf", "-j", inputFilename)
	cmd.Env = []string{"TESSDATA_PREFIX=/Users/bmcquee/git/openalpr/runtime_data/ocr/"}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logg.LogTo("OCR_TESSERACT", "cmd StdoutPipe failure:err:%v:", err)
		// log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		logg.LogTo("OCR_TESSERACT", "cmd Start failure:err:%v:", err)
		// log.Fatal(err)
	}

/*

	{
   "epoch_time" : 140658849846224,
   "processing_time_ms" : 328.659,
   "results" : [
      {
         "coordinates" : [
            {
               "y" : 277,
               "x" : 308
            },
            {
               "y" : 277,
               "x" : 643
            },
            {
               "y" : 432,
               "x" : 637
            },
            {
               "y" : 432,
               "x" : 303
            }
         ],
         "region" : "",
         "processing_time_ms" : 61.331001,
         "matches_template" : 0,
         "plate" : "EAZ6913",
         "confidence" : 88.822861,
         "region_confidence" : 0,
         "candidates" : [
            {
               "confidence" : 88.822861,
               "matches_template" : 0,
               "plate" : "EAZ6913"
            },
            {
               "confidence" : 71.701836,
               "matches_template" : 0,
               "plate" : "EA269T3"
            }
         ]
      }
   ]
}

*/

	type Candidate struct {
		Confidence float32
		Matches_Template int
		Plate string
	}
	type Coordinate struct {
		x int
		y int
	}
	type Result struct {
		Coordinates [4]Coordinate
		Region string
		Processing_Time_Ms  float32
		Matches_Template int
		Plate string
		Confidence float32
		Region_Confidence int
		Candidates [10]Candidate
	}
	var ocr_results struct {
		Epoch_Time int64
		Processing_Time_Ms  float32
		Results [1]Result
	}
	if err := json.NewDecoder(stdout).Decode(&ocr_results); err != nil {
		logg.LogTo("OCR_TESSERACT", "json NewDecoder Decode:failure:err:%v:", err)
		// log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		logg.LogTo("OCR_TESSERACT", "cmd Wait:failure:err:%v:", err)
		// log.Fatal(err)
	}
	fmt.Printf("plate:%v:confidence:%v:\n", (ocr_results.Results[0]).Plate, (ocr_results.Results[0]).Confidence)

	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	logg.LogTo("OCR_TESSERACT", "Error exec tesseract: %v %v", err, string(output))
	// 	return OcrResult{}, err
	// }

	// logg.LogTo("OCR_TESSERACT", "OCR output:%v:", string(output))
	// // get data from outfile
	// outBytes, err := ioutil.ReadFile(tmpOutFileName)
	// if err != nil {
	// 	logg.LogTo("OCR_TESSERACT", "Error getting data from out file: %v", err)
	// 	return OcrResult{}, err
	// }

	return OcrResult{
		Text: string((ocr_results.Results[0]).Plate),
	}, nil

}
