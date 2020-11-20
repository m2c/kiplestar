package httptool

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/m2c/kiplestar/commons"
	"github.com/m2c/kiplestar/commons/utils"
	"github.com/valyala/fasthttp"
)

type Client struct {
	// for config
	Host    string
	Port    int32
	Mode    string
	TimeOut time.Duration
	IsDebug bool
}

func (c *Client) initConfig() *Client {
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Mode == "" {
		c.Mode = "http"
	}
	if c.Port == 0 {
		if c.Mode == "https" {
			c.Port = 443
		} else {
			c.Port = 80
		}
	}
	if c.TimeOut == 0 {
		c.TimeOut = time.Second * 5
	}
	return c
}

// field in path must be surround whith '{' and '}', like 'id' in "http://xxx/{id}"
func (c *Client) parsePathParams(uri string, params map[string]string) (url string) {
	if uri != "" && uri[0] != '/' {
		uri = "/" + uri
	}
	url = fmt.Sprintf("%s://%s:%d%s", c.Mode, c.Host, c.Port, uri)
	if len(params) > 0 {
		for s, pm := range params {
			url = strings.ReplaceAll(url, "{"+s+"}", pm)
		}
	}
	return
}

func (c *Client) parseQueryParams(url string, params map[string]string) string {
	var urlAddress = ""
	if len(params) > 0 {
		for k, v := range params {
			if len(k) != 0 && len(v) != 0 {
				urlAddress = urlAddress + k + "=" + v + "&"
			}
		}
		urlAddress = strings.Trim(urlAddress, "&")
	}
	if urlAddress != "" {
		if strings.Contains(url, "?") {
			url = url + "&" + urlAddress
		} else {
			url = url + "?" + urlAddress
		}
	}
	return url
}

func (c *Client) parseHeaderParams(headerParams map[string]string, cookieParams map[string]string) (rs map[string]string) {
	rs = headerParams
	if len(cookieParams) > 0 {
		s := ""
		for n, v := range cookieParams {
			if s == "" {
				s = n + "=" + v
			} else {
				s = fmt.Sprintf("%s; %s=%s", s, n, v)
			}
		}
		rs["Cookie"] = s
	}

	return
}

func (c *Client) parseParams(uri string, params interface{}) (newUrl string, req interface{}, headers map[string]string, err error) {
	if params == nil {
		return
	}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	cookieParams := map[string]string{}

	pmType := reflect.TypeOf(params)
	if pmType.Kind() != reflect.Struct {
		err = errors.New("params need struct")
		return
	}
	pmValue := reflect.ValueOf(params)

	n := pmType.NumField()
	for i := 0; i < n; i++ {
		f := pmType.Field(i)
		tag := f.Tag
		in := tag.Get("in")
		if in == "" {
			err = errors.New("field has no tag 'in': " + f.Name)
			return
		}
		name := tag.Get("json")
		if name == "" {
			name = utils.ToLowerCamelCase(f.Name)
		}
		fv := pmValue.Field(i)
		switch in {
		case TAG_TYPE__BODY:
			if req != nil {
				err = errors.New("struct in 'body' should be only one.")
				return
			}
			req = fv.Interface()
		case TAG_TYPE__PATH:
			pathParams[name] = fv.String()
		case TAG_TYPE__QUERY:
			queryParams[name] = fv.String()
		case TAG_TYPE__HEADER:
			headerParams[name] = fv.String()
		case TAG_TYPE__COOKIE:
			cookieParams[name] = fv.String()
		default:
			err = errors.New("field tag has wrong 'in': " + in)
			return
		}
	}

	newUrl = c.parsePathParams(uri, pathParams)
	newUrl = c.parseQueryParams(newUrl, queryParams)
	headers = c.parseHeaderParams(headerParams, cookieParams)

	return
}

// Used by all method
// if has no headers, just put 3 params
func (c *Client) Request(method string, uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	c.initConfig()
	method = strings.ToUpper(method)

	var url string
	headerMap := map[string]string{}
	var newReq interface{}
	if req != nil {
		url, newReq, headerMap, err = c.parseParams(uri, req)
		if err != nil {
			return
		}
	} else {
		url = fmt.Sprintf("%s://%s:%d/%s", c.Mode, c.Host, c.Port, uri)
	}
	if _, ok := headerMap[fasthttp.HeaderContentType]; !ok {
		headerMap[fasthttp.HeaderContentType] = "application/json;charset=utf-8"
	}

	request := NewHttpRequest(url, newReq).SetMethod(method).SetTimeout(c.TimeOut).SetDebug(c.IsDebug)
	if len(headers) > 0 {
		for _, header := range headers {
			request.SetHeaders(header)
		}
	}
	if len(headerMap) > 0 {
		request.SetHeaders(headerMap)
	}

	result, err = request.Do()
	return
}

func (c *Client) ParseToResult(body []byte, res interface{}) (err error) {
	rs := Result{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return
	}

	// Unified handling of error codes
	if rs.Code != commons.OK {
		err = errors.New(fmt.Sprintf("An error occurred，Code: %d, Msg: %s, Time: %d", rs.Code, rs.Msg, rs.Time))
		return
	}

	if res != nil {
		rf := reflect.TypeOf(res)
		if rf.Kind() == reflect.Ptr {
			err = errors.New("Response need a ptr")
			return
		}
		rf = rf.Elem()
		if rf.Kind() == reflect.Struct || rf.Kind() == reflect.Slice || rf.Kind() == reflect.Array {
			err = json.Unmarshal(rs.Data, res)
			if err != nil {
				return
			}
		} else {
			rfv := reflect.ValueOf(res)
			if rfv.CanSet() {
				newRfv := reflect.ValueOf(rs)
				rfv.Set(newRfv)
			} else {
				err = errors.New("Response can not set value")
				return
			}
		}
	}

	return
}

func (c *Client) RequestAndParse(method string, uri string, req interface{}, resp interface{}, headers ...map[string]string) (err error) {
	var body []byte
	body, err = c.Request(method, uri, req, headers...)
	if err != nil {
		return
	}
	err = c.ParseToResult(body, resp)
	return
}
