package httptool

import "encoding/json"

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
