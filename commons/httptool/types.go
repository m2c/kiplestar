package httptool

import "encoding/json"

const (
	TAG_TYPE__QUERY  = "query"
	TAG_TYPE__BODY   = "body"
	TAG_TYPE__PATH   = "path"
	TAG_TYPE__HEADER = "header"
	TAG_TYPE__COOKIE = "cookie"
)

const (
	ContentTypeFormUrlencoded = "application/x-www-form-urlencoded"
	ContentTypeFormData       = "multipart/form-data"
	ContentTypeJson           = "application/json"
)

type RequestItem struct {
	Key   string
	Value string
}

type Result struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
	Time int64           `json:"time"`
}
