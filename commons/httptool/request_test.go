package httptool

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
)

type RequestTest struct {
	AdminId       int64  `json:"admin_id"`
	DeclineReason string `json:"decline_reason"`
	Id            int64  `json:"id"`
	Status        int64  `json:"status"`
}

func TestHttpRequest_Default(t *testing.T) {
	req := RequestTest{
		AdminId:       123,
		DeclineReason: "qwe",
		Id:            123,
		Status:        1,
	}
	body, err := NewHttpRequest(http.MethodGet, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit?ref_id=123", req).Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestHttpRequest_Form(t *testing.T) {
	req := RequestTest{
		AdminId:       123,
		DeclineReason: "qwe",
		Id:            123,
		Status:        1,
	}
	body, err := NewHttpRequest(http.MethodPost, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).SetHeaders(map[string]string{
		fasthttp.HeaderContentType: ContentTypeFormData,
	}).Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestHttpRequest_FormUrlencode(t *testing.T) {
	req := RequestTest{
		AdminId:       123,
		DeclineReason: "qwe",
		Id:            123,
		Status:        1,
	}
	body, err := NewHttpRequest(http.MethodPost, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).SetHeaders(map[string]string{
		fasthttp.HeaderContentType: ContentTypeFormUrlencoded,
	}).Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
