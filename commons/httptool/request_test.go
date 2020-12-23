package httptool

import (
	"net/http"
	"testing"
)

type RequestGormTest struct {
	AdminId       int64  `json:"admin_id"`
	DeclineReason string `json:"decline_reason"`
	Id            int64  `json:"id"`
	Status        int64  `json:"status"`
}

type RequestTest struct {
	AdminId []int64 `json:"admin_id"`
	Status  int64   `json:"status"`
}

func TestHttpRequest_Default(t *testing.T) {
	req := RequestTest{
		AdminId: []int64{1, 2, 3},
		Status:  1,
	}
	// url can take params also, struct 'req' will append to end of url.
	body, err := NewHttpRequest(http.MethodGet, "http://www.baidu.com?ref_id=123", req).Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestHttpRequest_Post(t *testing.T) {
	req := RequestTest{
		AdminId: []int64{1, 2, 3},
		Status:  1,
	}
	body, err := NewHttpRequest(http.MethodPost, "http://www.baidu.com", req).Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestHttpRequest_Form(t *testing.T) {
	req := RequestGormTest{
		AdminId:       123,
		DeclineReason: "qwe",
		Id:            123,
		Status:        1,
	}
	body, err := NewHttpRequest(http.MethodPost, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).RequestForm()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestHttpRequest_FormUrlencode(t *testing.T) {
	req := RequestGormTest{
		AdminId:       123,
		DeclineReason: "qwe",
		Id:            123,
		Status:        1,
	}
	body, err := NewHttpRequest(http.MethodPost, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).RequestFormUrlencoded()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
