package main

import (
	"encoding/json"
	"fmt"
	"github.com/m2c/kiplestar/commons"
	"github.com/m2c/kiplestar/commons/utils"
)

func main() {
	resp, err := utils.RiskInstance("https://www.baidu.com", "b8485198-9f51-4c00-8028-24722ab3bea5").
		Exec(utils.RiskLogin, utils.RiskLoginReq{AccountNo: "zhangkoe", IpAddress: "127.0.0.1"})
	println(err.Error())
	if resp != nil {
		println(fmt.Sprintf("%v", resp.IsBlocked()))
	}
}

func main1() {
	var t1 Test1
	t1.Mobile = "123123123123"
	println(utils.SensitiveStruct(t1))

	rsp := new(Resp)
	rsp.Account.Account = "koe"
	rsp.T.Mobile = "123123123123"
	rsp.Account.Pin.NewPin = "123333"
	bts, _ := json.Marshal(commons.BuildSuccess(rsp))
	println(utils.SensitiveFilter(string(bts)))
}

type Test struct {
	Mobile string `json:"mobile"`
	C
}

type Test1 struct {
	Mobile string `json:"mobile"`
}

type C struct {
	Account string `json:"account"`
	Pin     P      `json:"pin"`
}

type P struct {
	NewPin string `json:"new_pin"`
}

type Resp struct {
	Account C     `json:"account"`
	T       Test1 `json:"account1"`
	Pin     P     `json:"pin"`
}
