package utils

import (
	"encoding/json"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/valyala/fasthttp"
	"strings"
)

func Request(method string, url string, body interface{}, response interface{}) (err error) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(strings.ToUpper(method))
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
	if body != nil {
		binBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req.SetBody(binBody)
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		slog.Infof("Http Request Do Error %s", err.Error())
		return err
	}
	respBody := resp.Body()
	if response != nil {
		err = json.Unmarshal(respBody, response)
		if err != nil {
			return err
		}
	} else {
		slog.Info(string(respBody))
	}

	return nil

}
