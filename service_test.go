package kiplestar

import (
	"fmt"
	slog "kiplestar/commons/log"
	"testing"
)

func TestStart_Default_Server(t *testing.T) {
	slog.Info()
	server := GetKipleServerInstance()
	//http
	server.app.Default()
	err := server.StartServer(Mysql_service, Redis_service)
	if err != nil {
		fmt.Println(err.Error())
	}
}
