package utils

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	//out := map[string]interface{}{}
	err := Request(http.MethodGet, "http://www.baidu.com", nil, nil)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
