package httptool

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/m2c/kiplestar/commons"
	slog "github.com/m2c/kiplestar/commons/log"
	"go.uber.org/zap"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type HttpRequest struct {
	url     string
	method  string
	headers map[string]string
	params  interface{}
	Timeout time.Duration
	Logger  *zap.SugaredLogger
}

func NewHttpRequest(url string, params interface{}) *HttpRequest {
	return &HttpRequest{
		url:    url,
		params: params,
		method: fasthttp.MethodGet,
		headers: map[string]string{
			"Content-Length": "0",
			"Host":           url,
			"Accept":         "*/*",
			"Connection":     "keep-alive",
		},
	}
}

func (hr *HttpRequest) getBody() (body []byte, err error) {
	if hr.params == nil {
		return
	}
	switch strings.ToLower(hr.headers[fasthttp.HeaderContentType]) {
	case ContentTypeJson:
		body, err = json.Marshal(hr.params)
		if err != nil {
			return
		}
	case ContentTypeFormData:
		// TODO: support upload file later.
		req, e := FormatRequestParams(hr.params)
		if e != nil {
			err = e
			return
		}
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		defer writer.Close()
		for _, v := range req {
			writer.WriteField(v.Key, v.Value)
		}
		hr.headers[fasthttp.HeaderContentType] = writer.FormDataContentType()
		body = payload.Bytes()
	case ContentTypeFormUrlencoded:
		req, e := FormatRequestParams(hr.params)
		if e != nil {
			err = e
			return
		}
		payload := url.Values{}
		for _, v := range req {
			payload.Add(v.Key, v.Value)
		}
		body = []byte(payload.Encode())
	default:
		err = errors.New("content-type is not be supported.")
		return
	}
	return
}

func (hr *HttpRequest) initConfig() {
	if hr.method == "" {
		hr.method = fasthttp.MethodGet
	}
	if hr.Timeout == 0 {
		hr.Timeout = time.Second * 30
	}
	if _, ok := hr.headers[fasthttp.HeaderContentType]; !ok {
		hr.headers[fasthttp.HeaderContentType] = ContentTypeJson
	}
	if _, ok := hr.headers[fasthttp.HeaderUserAgent]; !ok {
		// set UserAgent, Avoid some server-side restrictions cannot be empty.
		hr.headers[fasthttp.HeaderUserAgent] = "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:22.0) Gecko/20100101 Firefox/22.0"
	}
	if _, ok := hr.headers[fasthttp.HeaderConnection]; !ok {
		hr.headers[fasthttp.HeaderConnection] = fasthttp.HeaderKeepAlive
	}
	if hr.Logger == nil && slog.Log != nil {
		hr.Logger = slog.Log
	}
}

//append request params to url
func (hr *HttpRequest) getTotalUrlByParams(urlStr string, params interface{}) (string, error) {
	ul, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	args, err := FormatRequestParams(params)
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		params := url.Values{}
		for _, v := range args {
			params.Add(v.Key, v.Value)
		}
		if ul.RawQuery == "" {
			ul.RawQuery = params.Encode()
		} else {
			ul.RawQuery += "&" + params.Encode()
		}
	}

	return ul.String(), nil
}

func (hr *HttpRequest) SetHeaders(headers map[string]string) *HttpRequest {
	for s, v := range headers {
		hr.headers[s] = v
	}
	return hr
}

func (hr *HttpRequest) SetTimeout(t time.Duration) *HttpRequest {
	hr.Timeout = t
	return hr
}

func (hr *HttpRequest) SetLogger(log *zap.SugaredLogger) *HttpRequest {
	hr.Logger = log
	return hr
}

func (hr *HttpRequest) SetMethod(method string) *HttpRequest {
	hr.method = method
	return hr
}

func (hr *HttpRequest) WithXRequestId(xid string) *HttpRequest {
	hr.SetHeaders(map[string]string{
		commons.X_REQUEST_ID: xid,
	})
	return hr
}

// method default: GET
// timeout default: 5 second, if using other timeout, please call function "SetTimeout" before
func (hr *HttpRequest) Do() (result []byte, err error) {
	if hr.url == "" {
		return nil, errors.New("url should not be empty")
	}
	hr.initConfig()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	if hr.method == fasthttp.MethodGet {
		hr.url, err = hr.getTotalUrlByParams(hr.url, hr.params)
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
		hr.headers[fasthttp.HeaderContentLength] = strconv.Itoa(len(body))
	}

	if len(hr.headers) > 0 {
		for s, v := range hr.headers {
			req.Header.Set(s, v)
		}
	}
	req.Header.SetMethod(hr.method)
	req.SetRequestURI(hr.url)

	if hr.Logger != nil {
		hr.Logger.Infof("[method]: %s [headers]: %#v [Body]: %#v", hr.method, req.Header.String(), string(req.Body()))
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if e := fasthttp.DoTimeout(req, resp, hr.Timeout); e != nil {
		err = e
		if hr.Logger != nil {
			hr.Logger.Infof("[url]: %s [error]: %s", hr.url, err.Error())
		}
		return
	}
	result = resp.Body()

	if hr.Logger != nil {
		hr.Logger.Infof("[url]: %s [response]: %s", hr.url, strings.Trim(string(result), "\n"))
	}

	return
}

func (hr *HttpRequest) Get() (result []byte, err error) {
	hr.method = fasthttp.MethodGet
	return hr.Do()
}

func (hr *HttpRequest) Put() (result []byte, err error) {
	hr.method = fasthttp.MethodPut
	return hr.Do()
}

func (hr *HttpRequest) Patch() (result []byte, err error) {
	hr.method = fasthttp.MethodPatch
	return hr.Do()
}

func (hr *HttpRequest) Delete() (result []byte, err error) {
	hr.method = fasthttp.MethodDelete
	return hr.Do()
}

func (hr *HttpRequest) Post() (result []byte, err error) {
	hr.method = fasthttp.MethodPost
	return hr.Do()
}

func (hr *HttpRequest) PostForm() (result []byte, err error) {
	hr.SetHeaders(map[string]string{
		fasthttp.HeaderContentType: ContentTypeFormData,
	})
	hr.method = fasthttp.MethodPost
	return hr.Do()
}

func (hr *HttpRequest) PostFormUrlencoded() (result []byte, err error) {
	hr.SetHeaders(map[string]string{
		fasthttp.HeaderContentType: ContentTypeFormUrlencoded,
	})
	hr.method = fasthttp.MethodPost
	return hr.Do()
}
