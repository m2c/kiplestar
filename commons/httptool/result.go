package httptool

import (
	"encoding/json"
	"github.com/m2c/kiplestar/commons"
)

type Result struct {
	Code commons.ResponseCode `json:"code"`
	Msg  string               `json:"msg"`
	Data json.RawMessage      `json:"data"`
	Time int64                `json:"time"`
}
