package commons

//define the error code
type ResponseCode int

const (
	// unknow error
	UnKnowError      ResponseCode = -1
	HttpRequestError ResponseCode = 1
	// ok
	OK             ResponseCode = 0
	ParameterError ResponseCode = 80000
	ValidateError  ResponseCode = 80001
	InternalError  ResponseCode = 80006
)

//global code and msg
var CodeMsg = map[ResponseCode]string{
	OK:             "success",
	UnKnowError:    "Internal server error, please try again later",
	ParameterError: "Parameter type miss match",
	ValidateError:  "Request Parameter has errors",
	InternalError:  "Server busy!",
}

//construct the code and msg
func GetCodeAndMsg(code ResponseCode) string {
	value, ok := CodeMsg[code]
	if ok {
		return value
	}
	return "{}"
}
