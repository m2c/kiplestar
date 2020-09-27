package httptool

import (
	"testing"
)

// every field need config tag 'in'
// field must be only one in 'body'
type ClientTest struct {
	Param          ClientTestBody `json:"Param" in:"body"`
	Debug          string         `json:"debug" in:"cookie"`
	ApplicationKey string         `json:"X-Application-Key" in:"header"`
	Authorization  string         `json:"Authorization" in:"header"`
	Token          string         `json:"token" in:"header"`
}

// body do not need tag 'in'
type ClientTestBody struct {
	AdminId       int64  `json:"admin_id"`
	DeclineReason string `json:"decline_reason"`
	Id            int64  `json:"id"`
	Status        int64  `json:"status"`
}

type ClientResp struct {
}

func TestClient_Request(t *testing.T) {
	client := Client{
		Host: "192.168.1.175",
		Port: 8080,
		Mode: "http",
	}

	req := ClientTest{
		Param: ClientTestBody{
			AdminId:       123,
			DeclineReason: "qwe",
			Id:            123,
			Status:        1,
		},
		Debug:          "abcxxx-77e1c83b-7bb0-437b-bc50-a7a58e5660ac",
		ApplicationKey: "5Syw49JVgmCDrGv5QBDbxtDvpTR2XxkF36Vr4EMVkvDVecJX",
		Authorization:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiMTNjMzJlZTEtZjIyZS00ZGQ5LTgxN2ItNjk3M2YzMTFhNGIzIiwibmFtZSI6IkRlc21vbmQgVW5kZXJ3b29kIiwiZW1haWwiOiJkZW1vQGRob21lLmlvIiwib2ZmaWNlX2VtYWlsIjpudWxsLCJwaG9uZV9jb3VudHJ5X2NvZGUiOiI2MCIsInBob25lIjoiNjAxMjAwMDExMTEiLCJhbHRfY291bnRyeV9jb2RlIjpudWxsLCJvZmZpY2VfcGhvbmUiOm51bGwsInBob3RvX3VybCI6Imh0dHBzOi8vczMtYXAtc291dGhlYXN0LTEuYW1hem9uYXdzLmNvbS9mbGlwYm94c2VjdXJlLWRob21lL3Byb2ZpbGVzL3Byb2ZpbGVfMTNjMzJlZTEtZjIyZS00ZGQ5LTgxN2ItNjk3M2YzMTFhNGIzXzEzM2FjMWM4MGNkMzRjOWJiMDUxZWRmODg1NGEwZTZkLmpwZyIsImNyZWF0ZWRfYXQiOiIyMDE3LTA2LTE5IDEzOjU4OjAzIiwidXBkYXRlZF9hdCI6IjIwMTgtMDktMDYgMDg6MjI6NTciLCJwaG90byI6ImMyYWE3N2JjLTM2MTctNDZmMi04NWVhLWI1Y2Y2YmJjYTYwNyIsImRlbGV0ZWRfYXQiOm51bGwsImZpcnN0X2xvZ2luX2F0IjoiMjAxOC0wOC0wMyAxNzo0NzowMSIsIm9sZF9waG9uZSI6bnVsbCwiaXNfZ3JvdXBfYWRtaW4iOmZhbHNlLCJpc19kZW1vX2FjY291bnQiOmZhbHNlLCJxcl9jb2RlIjpudWxsLCJxcmNvZGVfZXhwaXJlZF9hdCI6bnVsbCwicXJjb2RlX3VwZGF0ZWRfYXQiOm51bGwsInVzZXJfcHJvZmlsZV9pZCI6NTA3LCJzaG91bGRfbWlncmF0ZSI6IjAwMCIsInVtc190b2tlbl9zdGF0dXMiOiIwMDAiLCJtaWdyYXRpb25fc3RhdHVzIjoiMDAwIiwibWlncmF0ZWRfcGhvbmUiOm51bGwsImtiX3Rva2VuIjpudWxsLCJrYl9yZWZyZXNoX3Rva2VuIjpudWxsLCJpZGVudGl0aWVzIjpbeyJ1dWlkIjoiMGYzYzQyMDUtZGZmYS00ZjU1LTkzZDUtYzk4NWViYWFjNGEyIiwidHlwZSI6InVzZXJwYXNzIiwicHJvdmlkZXIiOiJhY2NvdW50cy5kaG9tZS5pbyIsInZlcmlmaWVkIjp0cnVlLCJpc19zZWxmX3JlZ2lzdGVyIjpmYWxzZSwiaWRlbnRpZmllciI6ImRlditkaG9tZTEzQGFwcGxhYi5teSIsIm9mZmljZV9pZGVudGlmaWVyIjpudWxsLCJyZXNldF9jb2RlIjpudWxsLCJhY3RpdmF0aW9uX2NvZGUiOiI3NTk1ZGFlOWNkZTgyMjE4MzM2YTU0NTdlZDlkNTVlYzg5OGM1MTYyM2Y3M2E2OWVlZmFhNTdhMmNjOTE5NGZjIiwiaXNfYWN0aXZhdGVkIjpmYWxzZSwiY3JlYXRlZF9hdCI6IjIwMTktMDItMTkgMDI6NTM6MDciLCJ1cGRhdGVkX2F0IjpudWxsLCJwaG9uZSI6bnVsbCwiZGVsZXRlZF9hdCI6bnVsbCwidXNlcl90eXBlIjoibmV3Iiwib2xkX3Bob25lIjpudWxsfSx7InV1aWQiOiI3ZmI1Y2IzNi1hNmExLTExZTgtOThkMC01MjkyNjlmYjE0NTkiLCJ0eXBlIjoidXNlcnBhc3MiLCJwcm92aWRlciI6ImFjY291bnRzLmRob21lLmlvIiwidmVyaWZpZWQiOnRydWUsImlzX3NlbGZfcmVnaXN0ZXIiOmZhbHNlLCJpZGVudGlmaWVyIjoiZGVtb0BkaG9tZS5pbyIsIm9mZmljZV9pZGVudGlmaWVyIjpudWxsLCJyZXNldF9jb2RlIjpudWxsLCJhY3RpdmF0aW9uX2NvZGUiOm51bGwsImlzX2FjdGl2YXRlZCI6dHJ1ZSwiY3JlYXRlZF9hdCI6IjIwMTgtMDgtMjMgMDY6NTY6NDUiLCJ1cGRhdGVkX2F0IjpudWxsLCJwaG9uZSI6bnVsbCwiZGVsZXRlZF9hdCI6bnVsbCwidXNlcl90eXBlIjoib2xkIiwib2xkX3Bob25lIjpudWxsfV0sInJvbGVzIjpbeyJ0eXBlIjoic3VwZXJhZG1pbiIsInJlc2lkZW5jZV91dWlkIjoiOTMwMjA3YjAtYmYxMC00MGJmLWE5ZTItMGIyOGJmZWNiZTFjIiwiY2FuX21hbmFnZV9yZXNpZGVuY2UiOm51bGx9XSwiaWF0IjoxNTk4NTc2ODg3LCJpc3MiOiJhY2NvdW50cy5kaG9tZS5pbyJ9.RnnUSkeI7RIjBDKhGkXdlK5rPFpV3aPalQaE9JPfnxE",
	}
	var body []byte
	var err error
	body, err = client.Request("POST", "/payment/web/v1.0/withdrawal/audit", req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	resp := ""
	err = client.ParseToResult(body, &resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
