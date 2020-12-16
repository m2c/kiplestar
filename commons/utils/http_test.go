package utils

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	//out := map[string]interface{}{}
	code, err := Request(http.MethodGet, "http://www.baidu.com", nil, nil, nil)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Println(code)
	}
}

func TestDoGetRequest(t *testing.T) {

	m1 := make(map[string]string)
	m1["a"] = "aa"
	m1["b"] = "bb"
	res, err := DoGetRequest("http://www.baidu.com", m1)
	fmt.Println(res)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
