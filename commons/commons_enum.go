package commons

import "github.com/kataras/iris/v12"

//define the error code
type ResponseCode int

var Application *iris.Application

const (
	// unknow error
	UnKnowError ResponseCode = -1
	// ok
	OK             ResponseCode = 0
	ParameterError ResponseCode = 80000
	ValidateError  ResponseCode = 80001
)

//global code and msg
var CodeMsg = map[ResponseCode]string{
	OK:             "success",
	UnKnowError:    "Internal server error, please try again later",
	ParameterError: "Parameter type miss match",
	ValidateError:  "Request Parameter has errors",
}

//construct the code and msg
func GetCodeAndMsg(code ResponseCode) string {
	value, ok := CodeMsg[code]
	if ok {
		return value
	}
	return ""
}
