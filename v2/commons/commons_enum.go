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
	TokenError     ResponseCode = 3
	CheckAuthError ResponseCode = 4
)

//global code and msg
var CodeMsg = map[ResponseCode]string{
	OK:             "success",
	UnKnowError:    "Internal server error, please try again later",
	HttpNotFound:   "Internal server error, Http not found",
	ParameterError: "Parameter type miss match",
	ValidateError:  "Request Parameter has errors",
	TokenError:     "Token Error",
	CheckAuthError: "Check Auth Error",
}

//construct the code and msg
func GetCodeAndMsg(code ResponseCode) string {
	value, ok := CodeMsg[code]
	if ok {
		return value
	}
	return "{}"
}

// msg will be used as default msg, and you can change msg with function 'BuildFailedWithMsg' or 'BuildSuccessWithMsg' or 'response.WithMsg' for once.
func RegisterCodeAndMsg(arr map[ResponseCode]string) {
	if len(arr) == 0 {
		return
	}
	for k, v := range arr {
		CodeMsg[k] = v
	}
}

const X_REQUEST_ID = "X-Request-Id"
