package httptool

import "testing"

type RequestTest struct {
	AdminId       int64  `json:"admin_id"`
	DeclineReason string `json:"decline_reason"`
	Id            int64  `json:"id"`
	Status        int64  `json:"status"`
}

func TestHttpRequest_Do(t *testing.T) {
	req := RequestTest{
		AdminId:       123,
		DeclineReason: "qwe",
		Id:            123,
		Status:        1,
	}
	body, err := NewHttpRequest("http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).SetMethod("POST").Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
