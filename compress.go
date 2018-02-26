package gotinypng

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"io/ioutil"

	"bytes"
	"fmt"
)

const (
	// tinypng api key
	CompressImgEmail = "YOUR_Email"
	CompressImgKey   = "YOUR_Key"
)

// @Title CompressImg
// @Description 图片压缩
// @Author lisea
// @Param path  	string 	"图片路径(支持网络与本地)"
// @Success compressImgContent([]byte) 	nil(error)
// @Failure nil([]byte) 	err(error)
func CompressImg(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, "https://api.tinify.com/shrink", nil)
	if err != nil {
		return nil, err
	}

	// source file for net or local url
	var reqBody []byte
	switch strings.Contains(path, "http") {
	case true: // net url
		req.Header.Set("Content-Type", "application/json")

		reqBody = bytes.NewBufferString(`{"source":{"url":"` + path + `"}}`).Bytes()
	default: // local path
		reqBody, err = ioutil.ReadFile(path)
	}

	if err != nil {
		return nil, err
	}

	// set body
	req.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	// Set Basic Auth
	req.SetBasicAuth(CompressImgEmail, CompressImgKey)

	fmt.Println(req.Header)

	c := &http.Client{
		Timeout: time.Second * 30, // set 3 second timeout
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 判断请求状态
	if res.StatusCode != http.StatusCreated {
		return nil, errors.New("HTTP " + res.Status)
	}

	// download compress file
	req, err = http.NewRequest(http.MethodGet, res.Header.Get("Location"), nil)
	if err != nil {
		return nil, err
	}
	// Set Basic Auth
	req.SetBasicAuth(CompressImgEmail, CompressImgKey)

	res, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 判断请求状态
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP " + res.Status)
	}

	return ioutil.ReadAll(res.Body)
}
