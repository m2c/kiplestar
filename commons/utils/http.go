package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/m2c/kiplestar/commons"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/valyala/fasthttp"
)

type BaseResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
	Time int64           `json:"time"`
}

type ProxyRequestHeader struct {
	ContentType string
}

var TimeOut = time.Second * 10

func ProxyRequest(method string, header http.Header, url string, body []byte) (response []byte, respHeader ProxyRequestHeader, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetBody(body)
	for s, v := range header {
		for _, v2 := range v {
			req.Header.Set(s, v2)
		}
	}
	req.Header.SetMethod(method)
	req.Header.Set(fasthttp.HeaderConnection, fasthttp.HeaderKeepAlive)
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.DoTimeout(req, resp, time.Second*5); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return nil, ProxyRequestHeader{}, err
	}

	return resp.Body(), ProxyRequestHeader{ContentType: string(resp.Header.ContentType())}, nil

}
func RequestFrom(method string, url string, body interface{}, response interface{}, header http.Header) (code int, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(strings.ToUpper(method))
	req.Header.Set(fasthttp.HeaderConnection, fasthttp.HeaderKeepAlive)
	req.SetRequestURI(url)
	if body != nil {
		tp := reflect.TypeOf(body)
		if tp.Kind() != reflect.Struct {
			return 0, errors.New("not struct")
		}
		values := req.PostArgs()
		ve := reflect.ValueOf(body)
		fieldNum := ve.NumField()
		for i := 0; i < fieldNum; i++ {
			if ve.Field(i).Type().Kind() == reflect.Struct {
				continue
			}
			values.Set(tp.Field(i).Tag.Get("json"), fmt.Sprintf("%v", ve.Field(i).Interface()))
		}
	}
	for s, v := range header {
		for _, v2 := range v {
			req.Header.Set(s, v2)
		}
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	slog.Infof("http request method : %s , url : %s , data : %v", method, url, body)
	if err := fasthttp.Do(req, resp); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return int(commons.UnKnowError), err
	}
	respBody := resp.Body()
	slog.Infof("http response method : %s , url : %s , body %s", string(respBody))
	if response != nil {
		baseResp := &BaseResponse{}
		err = json.Unmarshal(respBody, baseResp)
		if err != nil {
			return int(commons.ParameterError), err
		} else if baseResp.Code == 0 {
			if len(baseResp.Data) > 0 {
				err = json.Unmarshal(baseResp.Data, response)
				if err != nil {
					return int(commons.ParameterError), err
				}
			}
		} else {
			return baseResp.Code, fmt.Errorf("request do error %s", baseResp.Msg)
		}
	}
	return 0, nil

}
func Request(method string, url string, body interface{}, response interface{}, header http.Header) (code int, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(strings.ToUpper(method))
	req.Header.SetContentType("application/json")
	req.Header.Set(fasthttp.HeaderConnection, fasthttp.HeaderKeepAlive)
	req.SetRequestURI(url)
	if body != nil {
		binBody, err := json.Marshal(body)
		if err != nil {
			return int(commons.ParameterError), err
		}
		req.SetBody(binBody)
	}
	for s, v := range header {
		for _, v2 := range v {
			req.Header.Set(s, v2)
		}
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	slog.Infof("http request method : %s , url : %s , data : %v", method, url, body)
	if err := fasthttp.Do(req, resp); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return int(commons.UnKnowError), err
	}
	respBody := resp.Body()
	slog.Infof("http response method : %s , url : %s , body %s", string(respBody))
	if response != nil {
		baseResp := &BaseResponse{}
		err = json.Unmarshal(respBody, baseResp)
		if err != nil {
			return int(commons.ParameterError), err
		} else if baseResp.Code == 0 {
			if len(baseResp.Data) > 0 {
				err = json.Unmarshal(baseResp.Data, response)
				if err != nil {
					return int(commons.ParameterError), err
				}
			}
		} else {
			return baseResp.Code, fmt.Errorf(baseResp.Msg)
		}
	}
	return 0, nil

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

	slog.Infof("http request method : %s , url : %s", http.MethodGet, url)
	if err := fasthttp.DoTimeout(req, resp, TimeOut); err != nil {
		fmt.Println("Http Request Do Error %s" + err.Error())
		return "", err
	}
	respBody := resp.Body()
	slog.Infof("http response method : %s , url : %s , body %s", string(respBody))
	return string(respBody), nil

}

// DO post request
func DoPostRequest(url string, params map[string]string, header http.Header) (response string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	urlAddress := getRequestURL(url, params)
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(urlAddress)
	defer fasthttp.ReleaseResponse(resp)

	slog.Infof("http request method : %s , url : %s , data : %v", http.MethodPost, url, params)
	if err := fasthttp.DoTimeout(req, resp, TimeOut); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return "", err
	}
	respBody := resp.Body()
	slog.Infof("http response method : %s , url : %s , body %s", string(respBody))
	return string(respBody), nil
}

func DoPostJsonRequest(url string, params interface{}) (response string, err error) {
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

	slog.Infof("http request method : %s , url : %s , data : %v", http.MethodPost, url, params)
	if err := fasthttp.DoTimeout(req, resp, TimeOut); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return "", err
	}
	slog.Infof("http response method : %s , url : %s , body %s", string(resp.Body()))
	return string(resp.Body()), nil

}

// DO post request with customer header
func DoPostRequestWithHeader(url string, params map[string]string, header http.Header) (response string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	for s, v := range header {
		for _, v2 := range v {
			req.Header.Set(s, v2)
		}
	}
	req.Header.SetMethod("POST")
	urlAddress := getRequestURL(url, params)
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(urlAddress)
	defer fasthttp.ReleaseResponse(resp)

	slog.Infof("http request method : %s , url : %s , data : %v ", http.MethodPost, url, params)
	if err := fasthttp.DoTimeout(req, resp, TimeOut); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return "", err
	}
	respBody := resp.Body()
	slog.Infof("http response method : %s , url : %s , body %s ", string(respBody))
	return string(respBody), nil
}

// with header
func DoPostJsonRequestWithHeader(url string, params interface{}, header http.Header) (response string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	for s, v := range header {
		for _, v2 := range v {
			req.Header.Set(s, v2)
		}
	}

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

	slog.Infof("http request method : %s , url : %s , data : %v", http.MethodPost, url, params)
	if err := fasthttp.DoTimeout(req, resp, TimeOut); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return "", err
	}
	slog.Infof("http response method : %s , url : %s , body %s", string(resp.Body()))
	return string(resp.Body()), nil

}

func DoGetRequestWithHeader(url string, params map[string]string, header http.Header) (response string, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	for s, v := range header {
		for _, v2 := range v {
			req.Header.Set(s, v2)
		}
	}
	req.Header.SetMethod("GET")
	req.Header.SetContentType("application/json;charset=utf-8")
	urlAddress := getRequestURL(url, params)
	req.SetRequestURI(urlAddress)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	slog.Infof("http request method : %s , url : %s", http.MethodGet, url)
	if err := fasthttp.DoTimeout(req, resp, TimeOut); err != nil {
		fmt.Println("Http Request Do Error %s" + err.Error())
		return "", err
	}
	respBody := resp.Body()
	slog.Infof("http response method : %s , url : %s , body %s", string(respBody))
	return string(respBody), nil

}
