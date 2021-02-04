package httptool

import (
	"testing"
	"time"
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
	body, err := NewHttpRequest("http://127.0.0.1:8000?ref_id=123", req).Get()
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
	body, err := NewHttpRequest("http://127.0.0.1:8000", req).WithXRequestId("xxxx-xxxx-xxxx").Post()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestHttpRequest_PostWithHeaders(t *testing.T) {
	req := RequestTest{
		AdminId: []int64{1, 2, 3},
		Status:  1,
	}
	// default timeout is 30 * time.Second
	request := NewHttpRequest("http://127.0.0.1:8000", req).WithXRequestId("xxxx-xxxx-xxxx").SetTimeout(time.Second * 60)
	request.SetHeaders(map[string]string{
		"x-token": "xxxx-xxxx-xxxx",
	})
	body, err := request.Post()
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
	body, err := NewHttpRequest("http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).PostForm()
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
	body, err := NewHttpRequest("http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).PostFormUrlencoded()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
