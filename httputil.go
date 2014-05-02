package ocrworker

import (
	"io/ioutil"
	"net/http"
)

func getUrlContent(url string) (*[]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return &[]byte{}, err
	}

	defer resp.Body.Close()

	if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return &[]byte{}, err
	} else {
		return &bytes, nil
	}

}

func saveUrlContentToFileName(url, tmpFileName string) error {

	// TODO: current impl uses more memory than it needs to

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tmpFileName, bodyBytes, 0600)

}
