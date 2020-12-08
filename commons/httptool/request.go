package httptool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
	"strings"
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

func NewHttpRequest(method, url string, params interface{}) *HttpRequest {
	return &HttpRequest{
		Url:    url,
		Params: params,
		Method: strings.ToUpper(method),
		Headers: map[string]string{
			"Content-Length": "0",
			"Host":           url,
			"Accept":         "*/*",
			"Connection":     "keep-alive",
		},
		IsDebug: false,
	}
}

func (hr *HttpRequest) parseParams(params interface{}) (req map[string]string, err error) {
	if params == nil {
		return
	}
	tv := reflect.ValueOf(params)
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}
	tp := tv.Type()

	switch tp.Kind() {
	case reflect.Struct:
		req = map[string]string{}
		for i := 0; i < tp.NumField(); i++ {
			switch tp.Field(i).Type.Kind() {
			case reflect.String:
				req[hr.getJsonName(tp.Field(i))] = tv.Field(i).String()
			case reflect.Bool:
				req[hr.getJsonName(tp.Field(i))] = strconv.FormatBool(tv.Field(i).Bool())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				req[hr.getJsonName(tp.Field(i))] = fmt.Sprintf("%d", tv.Field(i).Int())
			case reflect.Float32, reflect.Float64:
				req[hr.getJsonName(tp.Field(i))] = fmt.Sprintf("%f", tv.Field(i).Float())
			default:
				err = fmt.Errorf("field [%s] type is invalid.", hr.getJsonName(tp.Field(i)))
				return
			}
		}
	case reflect.Map:
		v, ok := params.(map[string]string)
		if !ok {
			err = errors.New("params type is invalid, only struct and map is supported.")
			return
		}
		req = v
	default:
		err = errors.New("params type is invalid, only struct and map is supported.")
		return
	}
	return
}

//append request url
func (hr *HttpRequest) getRequestURL() (totalUrl string, err error) {
	ul, e := url.Parse(hr.Url)
	if e != nil {
		err = e
		return
	}
	args, e := hr.parseParams(hr.Params)
	if e != nil {
		err = e
		return
	}

	params := url.Values{}
	for k, v := range args {
		params.Add(k, v)
	}
	if ul.RawQuery == "" {
		ul.RawQuery = params.Encode()
	} else {
		ul.RawQuery += "&" + params.Encode()
	}

	totalUrl = ul.String()
	return
}

func (hr *HttpRequest) getJsonName(f reflect.StructField) string {
	name := f.Tag.Get("json")
	if name == "" {
		name = f.Name
	}
	return name
}

func (hr *HttpRequest) getBody() (body []byte, err error) {
	if hr.Params == nil {
		return
	}
	req, e := hr.parseParams(hr.Params)
	if e != nil {
		err = e
		return
	}
	switch strings.ToLower(hr.Headers[fasthttp.HeaderContentType]) {
	case ContentTypeJson:
		body, err = json.Marshal(hr.Params)
		if err != nil {
			return
		}
	case ContentTypeFormData:
		// TODO: support upload file later.
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		defer writer.Close()
		for k, v := range req {
			writer.WriteField(k, v)
		}
		hr.Headers[fasthttp.HeaderContentType] = writer.FormDataContentType()
		body = payload.Bytes()
	case ContentTypeFormUrlencoded:
		payload := url.Values{}
		for k, v := range req {
			payload.Add(k, v)
		}
		body = []byte(payload.Encode())
	default:
		err = errors.New("content-type is not be supported.")
		return
	}
	return
}

func (hr *HttpRequest) initConfig() {
	if hr.Method == "" {
		hr.Method = fasthttp.MethodGet
	}
	if hr.Timeout == 0 {
		hr.Timeout = time.Second * 5
	}
	if _, ok := hr.Headers[fasthttp.HeaderContentType]; !ok {
		hr.Headers[fasthttp.HeaderContentType] = ContentTypeJson
	}
	if _, ok := hr.Headers[fasthttp.HeaderUserAgent]; !ok {
		// set UserAgent, Avoid some server-side restrictions cannot be empty.
		hr.Headers[fasthttp.HeaderUserAgent] = "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:22.0) Gecko/20100101 Firefox/22.0"
	}
	if _, ok := hr.Headers[fasthttp.HeaderConnection]; !ok {
		hr.Headers[fasthttp.HeaderConnection] = fasthttp.HeaderKeepAlive
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
	hr.initConfig()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	if hr.Method == fasthttp.MethodGet {
		hr.Url, err = hr.getRequestURL()
		if err != nil {
			return
		}
	} else {
		var body []byte
		body, err = hr.getBody()
		if err != nil {
			return
		}
		req.SetBody(body)
		hr.Headers[fasthttp.HeaderContentLength] = strconv.Itoa(len(body))
	}

	if len(hr.Headers) > 0 {
		for s, v := range hr.Headers {
			req.Header.Set(s, v)
		}
	}
	req.Header.SetMethod(hr.Method)
	req.SetRequestURI(hr.Url)

	if hr.IsDebug {
		log.Printf("\033[1;32m\n [request]: %s\n[header]: %s\033[0m\n", string(req.Body()), req.Header.String())
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if e := fasthttp.DoTimeout(req, resp, hr.Timeout); e != nil {
		err = e
		if hr.IsDebug {
			log.Printf("\033[1;31m\n[url]: %s\n[error]: %s\033[0m\n", hr.Url, err.Error())
		}
		return
	}
	result = resp.Body()

	if hr.IsDebug {
		log.Printf("\033[1;32m\n[url]: %s\n[response]: %s\033[0m\n", hr.Url, string(result))
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
