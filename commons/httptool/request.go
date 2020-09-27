package httptool

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/valyala/fasthttp"
)

type HttpRequest struct {
	Url     string
	Method  string
	Headers map[string]string
	Timeout time.Duration
	Params  interface{}
}

func NewHttpRequest(url string, params interface{}) *HttpRequest {
	return &HttpRequest{
		Url:     url,
		Params:  params,
		Headers: map[string]string{},
	}
}

func (httpRequest *HttpRequest) SetHeaders(headers map[string]string) *HttpRequest {
	for s, v := range headers {
		httpRequest.Headers[s] = v
	}
	return httpRequest
}

func (httpRequest *HttpRequest) SetMethod(method string) *HttpRequest {
	httpRequest.Method = method
	return httpRequest
}

func (httpRequest *HttpRequest) SetTimeout(ts time.Duration) *HttpRequest {
	httpRequest.Timeout = ts
	return httpRequest
}

// method default: GET, if using other methods, please call function "SetMethod" before
// timeout default: 5 second, if using other timeout, please call function "SetTimeout" before
func (httpRequest *HttpRequest) Do() (result []byte, err error) {
	if httpRequest.Url == "" {
		return nil, errors.New("url should not be empty")
	}
	if httpRequest.Method == "" {
		httpRequest.Method = fasthttp.MethodGet
	}
	if httpRequest.Timeout == 0 {
		httpRequest.Timeout = time.Second * 5
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	if len(httpRequest.Headers) > 0 {
		for s, v := range httpRequest.Headers {
			req.Header.Set(s, v)
		}
	}
	req.Header.SetMethod(httpRequest.Method)
	if httpRequest.Method == fasthttp.MethodGet && httpRequest.Params != nil {
		err = httpRequest.getRequestURL()
		if err != nil {
			return
		}
	} else if httpRequest.Params != nil {
		body, err := json.Marshal(httpRequest.Params)
		if err != nil {
			return nil, err
		}
		req.SetBody(body)
	}
	req.SetRequestURI(httpRequest.Url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if e := fasthttp.DoTimeout(req, resp, httpRequest.Timeout); e != nil {
		err = e
		return
	}
	result = resp.Body()

	return
}

func (httpRequest *HttpRequest) Get() (result []byte, err error) {
	httpRequest.Method = fasthttp.MethodGet
	return httpRequest.Do()
}

func (httpRequest *HttpRequest) Post() (result []byte, err error) {
	httpRequest.Method = fasthttp.MethodPost
	return httpRequest.Do()
}

func (httpRequest *HttpRequest) Put() (result []byte, err error) {
	httpRequest.Method = fasthttp.MethodPut
	return httpRequest.Do()
}

func (httpRequest *HttpRequest) Patch() (result []byte, err error) {
	httpRequest.Method = fasthttp.MethodPatch
	return httpRequest.Do()
}

func (httpRequest *HttpRequest) Delete() (result []byte, err error) {
	httpRequest.Method = fasthttp.MethodDelete
	return httpRequest.Do()
}

//append request url
func (httpRequest *HttpRequest) getRequestURL() (err error) {
	params, ok := httpRequest.Params.(map[string]string)
	if !ok {
		return errors.New("'GET' request's param should be map")
	}

	url := httpRequest.Url
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
	httpRequest.Url = urlAddress

	return
}
