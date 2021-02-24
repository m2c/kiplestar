package commons

//define the error code
type ResponseCode int

const (
	// unknow error
	HttpNotFound ResponseCode = -2
	UnKnowError  ResponseCode = -1
	// ok
	OK             ResponseCode = 0
	ParameterError ResponseCode = 1
	ValidateError  ResponseCode = 2
)

//global code and msg
var CodeMsg = map[ResponseCode]string{
	OK:             "success",
	UnKnowError:    "Internal server error, please try again later",
	HttpNotFound:   "Internal server error, Http not found",
	ParameterError: "Parameter type miss match",
	ValidateError:  "Request Parameter has errors",
}

//construct the code and msg
func GetCodeAndMsg(code ResponseCode) string {
	value, ok := CodeMsg[code]
	if ok {
		return value
	}
	return "{}"
}

const X_REQUEST_ID = "X-Request-Id"
