package commons

import (
	"time"
)

type BaseResponse struct {
	Code ResponseCode `json:"code"`
	Msg  string       `json:"msg"`
	Data interface{}  `json:"data"`
	Time int64        `json:"time"`
}

type BaseResponseHeader struct {
	Code ResponseCode `json:"code"`
	Msg  string       `json:"msg"`
	Time int64        `json:"time"`
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

func BuildFailedHeader(code ResponseCode) BaseResponseHeader {
	return BaseResponseHeader{
		Code: code,
		Msg:  GetCodeAndMsg(code),
		Time: time.Now().UnixNano() / 1e6,
	}
}

func BuildSuccessHeader() BaseResponseHeader {
	msg := GetCodeAndMsg(OK)
	return BaseResponseHeader{
		Code: OK,
		Msg:  msg,
		Time: time.Now().UnixNano() / 1e6,
	}
}

func BuildWithHeader(header BaseResponseHeader, data interface{}) *BaseResponse {
	return &BaseResponse{
		Code: header.Code,
		Msg:  header.Msg,
		Data: data,
		Time: header.Time,
	}
}
