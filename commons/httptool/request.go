package httptool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"reflect"
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

func (hr *HttpRequest) getBody() (body []byte, err error) {
	switch hr.Headers["Content-Type"] {
	case ContentTypeJson:
		body, err = json.Marshal(hr.Params)
		if err != nil {
			return
		}
	case ContentTypeForm:
		tp := reflect.TypeOf(hr.Params)
		tv := reflect.ValueOf(hr.Params)
		switch tp.Kind() {
		case reflect.Struct:
			bt := bytes.Buffer{}
			flag := false
			for i := 0; i < tp.NumField(); i++ {
				if tv.Field(i).IsZero() {
					continue
				}
				if flag {
					bt.WriteByte('&')
				}
				bt.WriteString(tp.Field(i).Name)
				bt.WriteString("=")
				switch tp.Field(i).Type.Kind() {
				case reflect.String:
					bt.WriteString(tv.Field(i).String())
				case reflect.Bool:
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					bt.WriteString(fmt.Sprintf("%d", tv.Field(i).Int()))
				case reflect.Float32, reflect.Float64:
					bt.WriteString(fmt.Sprintf("%f", tv.Field(i).Float()))
				default:
					err = errors.New("field invalid type")
					return
				}
				flag = true
			}
			body = bt.Bytes()
		case reflect.Map:
			tv := reflect.ValueOf(hr.Params)
			bt := bytes.Buffer{}
			flag := false
			for _, key := range tv.MapKeys() {
				if tv.MapIndex(key).IsZero() {
					continue
				}
				if flag {
					bt.WriteByte('&')
				}
				bt.WriteString(key.String())
				bt.WriteString("=")
				bt.WriteString(tv.MapIndex(key).String())
				flag = true
			}
			body = bt.Bytes()
		default:
			err = errors.New("param is invalid: " + tp.Kind().String())
			return
		}
	}
	return
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

	if _, ok := hr.Headers[fasthttp.HeaderContentType]; !ok {
		hr.Headers[fasthttp.HeaderContentType] = ContentTypeJson
	}
	if _, ok := hr.Headers[fasthttp.HeaderUserAgent]; !ok {
		hr.Headers[fasthttp.HeaderUserAgent] = "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:22.0) Gecko/20100101 Firefox/22.0"
	}
	if hr.Method == fasthttp.MethodGet && hr.Params != nil {
		err = hr.getRequestURL()
		if err != nil {
			return
		}
	} else if hr.Params != nil {
		var body []byte
		body, err = hr.getBody()
		if err != nil {
			return
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
