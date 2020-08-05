package utils

import (
	"encoding/json"
	"fmt"
	"github.com/m2c/kiplestar/commons"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/valyala/fasthttp"
	"strings"
)

type BaseResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
	Time int64           `json:"time"`
}

func Request(method string, url string, body interface{}, response interface{}) (code int, err error) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(strings.ToUpper(method))
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	if body != nil {
		binBody, err := json.Marshal(body)
		if err != nil {
			return int(commons.ParameterError), err
		}
		req.SetBody(binBody)
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return int(commons.HttpRequestError), err
	}
	respBody := resp.Body()

	if response != nil {
		baseResp := &BaseResponse{}
		err = json.Unmarshal(respBody, baseResp)
		if err != nil {
			return int(commons.ParameterError), err
		} else if baseResp.Code == 0 && len(baseResp.Data) > 0 {
			err = json.Unmarshal(baseResp.Data, response)
			if err != nil {
				return int(commons.ParameterError), err
			}
		} else {
			return baseResp.Code, fmt.Errorf("Request do error %s", baseResp.Msg)
		}
	} else {
		slog.Info(string(respBody))
	}

	return 0, nil

}
