package kiplestar

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"testing"
)

func TestStart_Default_Server(t *testing.T) {
	server := GetKipleServerInstance()
	//http
	server.app.Default()
	err := server.StartServer(Mysql_service)
	if err != nil {
		fmt.Println(err.Error())
	}
	server.WaitClose(iris.WithoutBodyConsumptionOnUnmarshal)
}
