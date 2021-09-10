package main

import (
	"encoding/json"
	"github.com/m2c/kiplestar/commons"
	"github.com/m2c/kiplestar/commons/utils"
)

func main() {
	var t1 Test1
	t1.Mobile = 11
	println(utils.SensitiveStruct(nil))

	rsp := new(Resp)
	rsp.Account = "koe"
	bts,_ := json.Marshal(commons.BuildSuccess(rsp))
	println(utils.SensitiveFilter(string(bts)))
}

type Test struct {
	Mobile string `json:"mobile"`
 	C
}

type Test1 struct {
	Mobile int `json:"mobile"`
}

type C struct {
	Account string `json:"account"`
}

type Resp struct {
	Account string `json:"account"`
	Account1 string `json:"account1"`
	Account2 string `json:"account2"`
	Account3 string `json:"account3"`
	Account4 string `json:"account4"`
	Account5 string `json:"account5"`
	Account6 string `json:"account6"`
	Account7 string `json:"account7"`
	Account8 string `json:"account8"`
	Account9 string `json:"account9"`
	Account10 string `json:"account10"`
}

