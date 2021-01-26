# httptool

### base request

A base http tool for calling API.

- demo for a simple get http

```golang
req := map[string]string{
    "name": "test",
}
body, err := NewHttpRequest(http.MethodGet, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).Do()
if err != nil {
   // handle err
}
// handle body
```

or

```golang
type TestRequest struct {
    Name string `json:"name"`
}
req := TestRequest{ Name: "test" }
body, err := NewHttpRequest(http.MethodGet, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit", req).Do()
if err != nil {
   // handle err
}
// handle body
```

- demo for a post http

```golang
req := RequestTest{
    AdminId:       123,
    DeclineReason: "qwe",
    Id:            123,
    Status:        1,
}
// if has query params, must put behind of the url.
body, err := NewHttpRequest(http.MethodPost, "http://192.168.1.175:8080/payment/web/v1.0/withdrawal/audit?name=test", req).Do()
if err != nil {
    // handle err
}
// handle body
```

- more examples, please see request_test.go
