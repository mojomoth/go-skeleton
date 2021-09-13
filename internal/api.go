package internal

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	API = "http://3.36.180.31:11337"
)

func Call(method string, path string, body []byte) ([]byte, int, error) {
	req, _ := http.NewRequest(method, API+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	res, _ := ioutil.ReadAll(resp.Body)
	return res, resp.StatusCode, nil
}
