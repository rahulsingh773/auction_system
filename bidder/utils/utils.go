package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"
)

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

func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := make([]byte, length)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func GetRandomNumber() float32 {
	rand.Seed(time.Now().UnixNano())
	num := rand.Float32() * 1000
	return num
}
