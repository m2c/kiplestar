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
const X_SPAN_ID = "X-Span-Id"

const (
	//HTTP CODE
	HTTP_200 int = 200
	HTTP_400 int = 400
	HTTP_FAIL_1 int = -1
	//LOG TYPE
	LOG_TYPE_CURL string = "CURL"
	LOG_TYPE_MYSQL string = "MYSQL"
	LOG_TYPE_PUSH string = "PUSH"
	LOG_TYPE_MQ string = "MQ"
	LOG_TYPE_REDIS string = "REDIS"
	//LOG FIELDS
	LOG_FIELD_TYPE string = "type"
	LOG_FIELD_URL string = "url"
	LOG_FIELD_METHOD string = "method"
	LOG_FIELD_INPUT string = "input"
	LOG_FIELD_RSP_TIME string = "rspTime"
	LOG_FIELD_MESSAGE string = "message"
	LOG_FIELD_RESPONSE string = "response"
	LOG_FIELD_HTTP_CODE string = "httpCode"
	LOG_FIELD_BUSINESS_CODE string = "businessCode"
	LOG_FIELD_SERVICE string = "service"
	LOG_FIELD_HEADER string = "header"
	LOG_FIELD_ARGS string = "args"
	LOG_FIELD_CLIENT string = "client"
)