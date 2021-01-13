package db_log

import (
	"github.com/kataras/iris/v12"
	"github.com/m2c/kiplestar/commons/utils"
	"log"
	"os"
)

type DbLogger struct {
	log        *log.Logger
	XRequestId string
}

func NewDbLogger(ctx iris.Context) *DbLogger {
	l := log.New(os.Stdout, "Db Logger:\t", log.LstdFlags|log.Lshortfile)
	return &DbLogger{
		log:        l,
		XRequestId: utils.GetXRequestID(ctx),
	}
}

func (l *DbLogger) Printf(format string, v ...interface{}) {
	l.log.Printf("X-Request-Id: "+l.XRequestId+"\t"+format, v...)
}

func (l *DbLogger) Print(v ...interface{}) {
	l.log.Print("X-Request-Id: "+l.XRequestId+"\t", v)
}

func (l *DbLogger) Println(v ...interface{}) {
	l.log.Println("X-Request-Id: "+l.XRequestId+"\t", v)
}
