package reponse

import (
	"errors"
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

// args can be empty. Or the first arg should be 'ResponseCode', and second arg should be 'Data'.
func NewCommonResponse(args ...interface{}) *CommonResponse {
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

func (c *CommonResponse) WithCode(code commons.ResponseCode) *CommonResponse {
	c.Code = code
	c.Msg = commons.GetCodeAndMsg(code)
	return c
}

func (c *CommonResponse) WithMsg(msg string) *CommonResponse {
	c.Msg = msg
	return c
}

func (c *CommonResponse) WithTraceId(tid string) *CommonResponse {
	c.TraceId = tid
	return c
}
