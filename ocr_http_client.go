package ocrworker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	ApiEndpointUrl string // the url of the server, eg, http://api.openocr.net
}

func NewHttpClient(apiEndpointUrl string) *HttpClient {
	return &HttpClient{
		ApiEndpointUrl: apiEndpointUrl,
	}
}

func (c HttpClient) DecodeImageUrl(u string, e OcrEngineType) (string, error) {

	ocrRequest := OcrRequest{
		ImgUrl:     u,
		EngineType: e,
	}

	// create JSON for POST reqeust
	jsonBytes, err := json.Marshal(ocrRequest)
	if err != nil {
		return "", err
	}

	// create a client
	client := &http.Client{}

	// create POST request
	apiUrl := c.OcrApiEndpointUrl()
	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}

	// send POST request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Got error status response: %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBytes), nil

}

// Get the url of the OCR API endpoint, eg, http://api.openocr.net/ocr
func (c HttpClient) OcrApiEndpointUrl() string {
	return fmt.Sprintf("%v/%v", c.ApiEndpointUrl, "ocr")
}
