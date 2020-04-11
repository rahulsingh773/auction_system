package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

//common method to send http request
func SendRestRequest(method string, path string, data string) ([]byte, int, error) {
	req, err := http.NewRequest(method, path, bytes.NewBufferString(data))
	if nil != err {
		log.Printf("SendRestRequest: failed to create request: %v", err)
		return nil, 0, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	req_dump, _ := httputil.DumpRequest(req, true)
	log.Printf("Sending Request: %v", string(req_dump))

	c := &http.Client{}
	resp, err := c.Do(req)
	if nil != err {
		log.Printf("SendRestRequest: %v() failed: %v", method, err)
		return nil, 0, err
	}

	res_dump, _ := httputil.DumpResponse(resp, true)
	log.Printf("Received Response: %v", string(res_dump))

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Printf("SendRestRequest: ReadAll() failed: %v", err)
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
