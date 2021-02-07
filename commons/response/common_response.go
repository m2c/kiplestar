package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/m2c/kiplestar/commons"
	"time"
)

type CommonResponse struct {
	Code    commons.ResponseCode `json:"code"`
	Msg     string               `json:"msg,omitempty"`
	Data    interface{}          `json:"data,omitempty"`
	TraceId string               `json:"trace_id,omitempty"`
	Time    int64                `json:"time,omitempty"`
}

type HttpResponse struct {
	Code commons.ResponseCode `json:"code"`
	Msg  string               `json:"msg,omitempty"`
	Data json.RawMessage      `json:"data,omitempty"`
	Time int64                `json:"time,omitempty"`
}

// args can be empty. Or the first arg should be 'ResponseCode', and second arg should be 'Data'.
func NewResponse(args ...interface{}) *CommonResponse {
	res := &CommonResponse{
		Time: time.Now().UnixNano() / 1e6,
	}
	switch len(args) {
	case 0:
		res.Code = commons.OK
	case 1:
		code, ok := args[0].(commons.ResponseCode)
		if !ok {
			panic(errors.New("code type invalid"))
		}
		res.Code = code
		res.Msg = commons.GetCodeAndMsg(code)
	case 2:
		code, ok := args[0].(commons.ResponseCode)
		if !ok {
			panic(errors.New("code type invalid"))
		}
		res.Code = code
		res.Data = args[1]
	default:
		panic(errors.New("the number of args is error"))
	}
	return res
}

func (c *CommonResponse) WithMsg(msg string) *CommonResponse {
	c.Msg = msg
	return c
}

func (c *CommonResponse) WithTraceId(tid string) *CommonResponse {
	c.TraceId = tid
	return c
}

func ParseResponse(body []byte, resp interface{}) error {
	rs := HttpResponse{}
	err := json.Unmarshal(body, &rs)
	if err != nil {
		return err
	}
	if rs.Code != commons.OK {
		return fmt.Errorf("request get a faid response - %#v", rs)
	}
	if resp == nil {
		return errors.New("response can not parse to nil address")
	}
	err = json.Unmarshal(rs.Data, resp)
	if err != nil {
		return err
	}
	return nil
}