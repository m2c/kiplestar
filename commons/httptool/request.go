package httptool

import (
	"encoding/json"
	"errors"
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"time"

	"github.com/valyala/fasthttp"
)

type HttpRequest struct {
	Url     string
	Method  string
	Headers map[string]string
	Timeout time.Duration
	Params  interface{}
	IsDebug bool
}

func NewHttpRequest(url string, params interface{}) *HttpRequest {
	return &HttpRequest{
		Url:    url,
		Params: params,
		Headers: map[string]string{
			"Content-Length": "0",
			"Host":           url,
			"Accept":         "*/*",
			"Connection":     "keep-alive",
		},
		IsDebug: false,
	}
}

func (hr *HttpRequest) SetHeaders(headers map[string]string) *HttpRequest {
	for s, v := range headers {
		hr.Headers[s] = v
	}
	return hr
}

func (hr *HttpRequest) SetDebug(debug bool) *HttpRequest {
	hr.IsDebug = debug
	return hr
}

func (hr *HttpRequest) SetMethod(method string) *HttpRequest {
	hr.Method = method
	return hr
}

func (hr *HttpRequest) SetTimeout(ts time.Duration) *HttpRequest {
	hr.Timeout = ts
	return hr
}

// method default: GET, if using other methods, please call function "SetMethod" before
// timeout default: 5 second, if using other timeout, please call function "SetTimeout" before
func (hr *HttpRequest) Do() (result []byte, err error) {
	if hr.Url == "" {
		return nil, errors.New("url should not be empty")
	}
	if hr.Method == "" {
		hr.Method = fasthttp.MethodGet
	}
	if hr.Timeout == 0 {
		hr.Timeout = time.Second * 5
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	if _, ok := hr.Headers["Content-Type"]; !ok {
		hr.Headers["Content-Type"] = "application/json"
	}
	if hr.Method == fasthttp.MethodGet && hr.Params != nil {
		err = hr.getRequestURL()
		if err != nil {
			return
		}
	} else if hr.Params != nil {
		body, err := json.Marshal(hr.Params)
		if err != nil {
			return nil, err
		}
		req.SetBody(body)
		hr.Headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	}

	if len(hr.Headers) > 0 {
		for s, v := range hr.Headers {
			req.Header.Set(s, v)
		}
	}
	req.Header.SetMethod(hr.Method)
	req.SetRequestURI(hr.Url)

	if hr.IsDebug {
		slog.Debugf("\033[1;32m[url]: %s\n [request]: %s\n [header]: %s\033[0m\n", hr.Url, string(req.Body()), req.Header.String())
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if e := fasthttp.DoTimeout(req, resp, hr.Timeout); e != nil {
		err = e
		if hr.IsDebug {
			slog.Debugf("\033[1;31m[url]: %s\n [error]: %s\033[0m\n", hr.Url, err.Error())
		}
		return
	}
	result = resp.Body()

	if hr.IsDebug {
		slog.Debugf("\033[1;32m[url]: %s\n [response]: %s\033[0m\n", hr.Url, string(result))
	}

	return
}

func (hr *HttpRequest) Get() (result []byte, err error) {
	hr.Method = fasthttp.MethodGet
	return hr.Do()
}

func (hr *HttpRequest) Post() (result []byte, err error) {
	hr.Method = fasthttp.MethodPost
	return hr.Do()
}

func (hr *HttpRequest) Put() (result []byte, err error) {
	hr.Method = fasthttp.MethodPut
	return hr.Do()
}

func (hr *HttpRequest) Patch() (result []byte, err error) {
	hr.Method = fasthttp.MethodPatch
	return hr.Do()
}

func (hr *HttpRequest) Delete() (result []byte, err error) {
	hr.Method = fasthttp.MethodDelete
	return hr.Do()
}

//append request url
func (hr *HttpRequest) getRequestURL() (err error) {
	params, ok := hr.Params.(map[string]string)
	if !ok {
		return errors.New("'GET' request's param should be map")
	}

	url := hr.Url
	var urlAddress string
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
	hr.Url = urlAddress

	return
}
