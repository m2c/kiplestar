package httptool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/m2c/kiplestar/commons"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	// for config
	Host       string
	Port       int32
	Mode       string
	TimeOut    time.Duration
	IsPrintLog bool
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

// first args for query, second args for path.
func (c *Client) getTotalUrl(uri string, args ...map[string]string) (u string, err error) {
	var ul *url.URL
	if strings.HasPrefix(uri, "http") {
		ul, err = url.Parse(uri)
		if err != nil {
			return
		}
	} else {
		ul, err = url.Parse(fmt.Sprintf("%s://%s:%d/%s", c.Mode, c.Host, c.Port, uri))
		if err != nil {
			return
		}
	}
	switch len(args) {
	case 0:
		u = ul.String()
	case 1:
		params := url.Values{}
		for k, v := range args[0] {
			params.Add(k, v)
		}
		ul.RawQuery = params.Encode()
	case 2:
		// use json tag of one field, and field in path must be surround whith '{' and '}', like 'id' in "http://xxx/{id}"
		uri = ul.Path
		for s, pm := range args[1] {
			uri = strings.ReplaceAll(uri, "{"+s+"}", pm)
		}
		ul.Path = uri

		params := url.Values{}
		for k, v := range args[0] {
			params.Add(k, v)
		}
		ul.RawQuery = params.Encode()
	default:
		err = errors.New("args max number is 2.")
		return
	}
	u = ul.String()
	return
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

func (c *Client) parseParamsWithTag(params interface{}) (req interface{}, pathMap, queryMap, headers map[string]string, err error) {
	headers = map[string]string{}
	if params == nil {
		return
	}
	pmType := reflect.TypeOf(params)
	if pmType.Kind() == reflect.Ptr {
		pmType = pmType.Elem()
	}
	switch pmType.Kind() {
	case reflect.Struct:
		pathParams := map[string]string{}
		queryParams := map[string]string{}
		headerParams := map[string]string{}
		cookieParams := map[string]string{}
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
				name = ToLowerCamelCase(f.Name)
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
		pathMap = pathParams
		queryMap = queryParams
		headers = c.parseHeaderParams(headerParams, cookieParams)
	case reflect.Map:
		req = params
	default:
		err = errors.New("params type is invalid, only struct and map is supported.")
		return
	}
	return
}

// parse the whole host with mode and prot to config, like 'http://127.0.0.1:8000'
func (c *Client) ParseDomain(totalHost string) error {
	if !strings.HasPrefix(totalHost, "http") {
		return errors.New("host is invalid")
	}
	ul, err := url.Parse(totalHost)
	if err != nil {
		return err
	}
	i := strings.Index(ul.Host, ":")
	if i == -1 {
		c.Host = ul.Host
	} else {
		c.Host = ul.Host[:i]
		port, err := strconv.ParseInt(ul.Host[i+1:], 10, 32)
		if err != nil {
			return err
		}
		c.Port = int32(port)
	}
	c.Mode = ul.Scheme
	return nil
}

// return the whole host with mode and port.
func (c *Client) GetDomain() string {
	return fmt.Sprintf("%s://%s:%d", c.Mode, c.Host, c.Port)
}

// Used by all method.
// fields need put 'in' into tag of req, or use map[string]string as req interface{}.
// tag 'in' supports 'header,path,query,cookie,body'.
func (c *Client) RequestWithAllTypeParams(method string, uri string, req interface{}) (result []byte, err error) {
	c.initConfig()
	method = strings.ToUpper(method)

	var urlStr string
	var pathMap, queryMap, headers map[string]string
	var newReq interface{}
	if req != nil {
		newReq, pathMap, queryMap, headers, err = c.parseParamsWithTag(req)
		if err != nil {
			return
		}
		urlStr, err = c.getTotalUrl(uri, queryMap, pathMap)
		if err != nil {
			return
		}
	} else {
		urlStr, err = c.getTotalUrl(uri)
		if err != nil {
			return
		}
	}

	request := NewHttpRequest(method, urlStr, newReq).SetTimeout(c.TimeOut).SetPrintLog(c.IsPrintLog)
	if len(headers) > 0 {
		request.SetHeaders(headers)
	}

	result, err = request.Do()
	return
}

// Not recommended to request with a url which contain host, host should controlled by *Client.
// although this func can handle the url with "http://".
func (c *Client) Request(method string, uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	c.initConfig()
	uri, err = c.getTotalUrl(uri)
	if err != nil {
		return
	}
	request := NewHttpRequest(method, uri, req).SetTimeout(c.TimeOut).SetPrintLog(c.IsPrintLog)
	if len(headers) > 0 {
		for _, hm := range headers {
			request.SetHeaders(hm)
		}
	}

	result, err = request.Do()
	return
}

func (c *Client) ParseCommonResponse(body []byte, resp interface{}) (err error) {
	rs := Result{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return
	}

	// Unified handling of error codes
	if rs.Code != int(commons.OK) {
		err = errors.New(fmt.Sprintf("An error occurred，Code: %d, Msg: %s, Time: %d", rs.Code, rs.Msg, rs.Time))
		return
	}

	if resp == nil {
		err = errors.New("Response can not parse nil address")
		return
	}
	rf := reflect.TypeOf(resp)
	if rf.Kind() != reflect.Ptr {
		err = errors.New("Response need a ptr")
		return
	}
	rf = rf.Elem()
	if rf.Kind() == reflect.Struct || rf.Kind() == reflect.Slice || rf.Kind() == reflect.Array {
		err = json.Unmarshal(rs.Data, resp)
		if err != nil {
			return
		}
	} else {
		rfv := reflect.ValueOf(resp)
		if rfv.CanSet() {
			newRfv := reflect.ValueOf(rs)
			rfv.Set(newRfv)
		} else {
			err = errors.New("Response can not set value")
			return
		}
	}
	return
}

// if req is not nil, this func will parse req and append to uri.
func (c *Client) Get(uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	return c.Request(http.MethodGet, uri, req, headers...)
}

func (c *Client) Post(uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	return c.Request(http.MethodPost, uri, req, headers...)
}

func (c *Client) PostForm(uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	headers = append(headers, map[string]string{
		fasthttp.HeaderContentType: ContentTypeFormData,
	})
	return c.Request(http.MethodPost, uri, req, headers...)
}

func (c *Client) PostFormUrlencoded(uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	headers = append(headers, map[string]string{
		fasthttp.HeaderContentType: ContentTypeFormUrlencoded,
	})
	return c.Request(http.MethodPost, uri, req, headers...)
}

func (c *Client) Put(uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	return c.Request(http.MethodPut, uri, req, headers...)
}

func (c *Client) Patch(uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	return c.Request(http.MethodPatch, uri, req, headers...)
}

func (c *Client) Delete(uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	return c.Request(http.MethodDelete, uri, req, headers...)
}
