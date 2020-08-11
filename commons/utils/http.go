package utils

import (
	"encoding/json"
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type BaseResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
	Time int64           `json:"time"`
}

func Request(method string, url string, body interface{}, response interface{}) (err error) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(strings.ToUpper(method))
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	if body != nil {
		binBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req.SetBody(binBody)
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return err
	}
	respBody := resp.Body()

	if response != nil {
		baseResp := &BaseResponse{}
		err = json.Unmarshal(respBody, baseResp)
		if err != nil {
			return err
		} else if len(baseResp.Data) > 0 {
			err = json.Unmarshal(baseResp.Data, baseResp)
			if err != nil {
				return err
			}
		}
	} else {
		slog.Info(string(respBody))
	}

	return nil

}

//append request url
func getRequestURL(url string, params map[string]string) string {
	var urlAddress = ""
	lastCharctor := url[len(url)-1:]
	if lastCharctor == "?" {
		urlAddress = url + urlAddress
	} else {
		urlAddress = url + "?" + urlAddress
	}
	for k, v := range params {
		if len(k) != 0 && len(v) != 0 {
			urlAddress = urlAddress + k + "=" + v + "&"
		}
	}
	return urlAddress
}

func DoGetRequest(url string, params map[string]string) (response string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json;charset=utf-8")
	urlAddress := getRequestURL(url, params)
	req.SetRequestURI(urlAddress)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.DoTimeout(req, resp, time.Second*10); err != nil {
		fmt.Println("Http Request Do Error %s" + err.Error())
		return "", err
	}
	respBody := resp.Body()
	return string(respBody), nil

}

// DO post request
func DoPostRequest(url string, params map[string]string) (response string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	urlAddress := getRequestURL(url, params)
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(urlAddress)
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.DoTimeout(req, resp, time.Second*10); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return "", err
	}
	respBody := resp.Body()
	return string(respBody), nil
}

func DoPostJsonRequest(url string, params map[string]string) (response string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json;charset=utf-8")
	req.SetRequestURI(url)
	if params != nil {
		binBody, err := json.Marshal(params)
		if err != nil {
			return "", err
		}
		req.SetBody(binBody)
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.DoTimeout(req, resp, time.Second*10); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return "", err
	}
	return string(resp.Body()), nil

}
