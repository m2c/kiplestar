package commons

import (
	"time"
)

type BaseResponse struct {
	Code ResponseCode `json:"code"`
	Msg  string       `json:"msg,omitempty"`
	Data interface{}  `json:"data,omitempty"`
	Time int64        `json:"time,omitempty"`
}

// return struct of the response code and msg
func BuildResponse(code ResponseCode, msg string, data interface{}) *BaseResponse {
	return &BaseResponse{code, msg, data, time.Now().UnixNano() / 1e6}
}

func BuildSuccess(data interface{}) *BaseResponse {

	return &BaseResponse{Code: OK, Msg: GetCodeAndMsg(OK), Data: data, Time: time.Now().UnixNano() / 1e6}
}
func BuildSuccessWithMsg(msg string, data interface{}) *BaseResponse {

	return &BaseResponse{Code: OK, Msg: msg, Data: data, Time: time.Now().UnixNano() / 1e6}
}
func BuildSuccessWithCode(code ResponseCode, data interface{}) *BaseResponse {
	msg := GetCodeAndMsg(code)
	return &BaseResponse{Code: code, Msg: msg, Data: data, Time: time.Now().UnixNano() / 1e6}
}
func BuildSuccessWithNoData(code ResponseCode) *BaseResponse {
	msg := GetCodeAndMsg(code)
	return &BaseResponse{Code: code, Msg: msg, Time: time.Now().UnixNano() / 1e6}
}
func BuildFailed(code ResponseCode) *BaseResponse {
	return &BaseResponse{
		Code: code,
		Msg:  GetCodeAndMsg(code),
		Data: struct{}{},
		Time: time.Now().UnixNano() / 1e6,
	}
}
func BuildFailedWithMsg(code ResponseCode, msg string) *BaseResponse {
	message := msg
	if len(msg) == 0 {
		message = GetCodeAndMsg(code)
	}
	return &BaseResponse{
		Code: code,
		Msg:  message,
		Data: struct{}{},
		Time: time.Now().UnixNano() / 1e6,
	}
}
func BuildFailedWithCode(code ResponseCode, data interface{}) *BaseResponse {
	msg := GetCodeAndMsg(code)
	if data == nil {
		data = struct{}{}
	}
	return &BaseResponse{
		Code: code,
		Msg:  msg,
		Data: data,
		Time: time.Now().UnixNano() / 1e6,
	}
}
