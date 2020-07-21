package middleware

import (
	"fmt"
	"github.com/kataras/iris/v12"
	slog "github.com/m2c/kiplestar/commons/log"
	"time"
)

func LoggerHandler(ctx iris.Context) {
	p := ctx.Request().URL.Path
	method := ctx.Request().Method
	start := time.Now().UnixNano() / 1e6
	ip := ctx.Request().RemoteAddr
	ctx.Request().URL.String()
	ctx.Request().UserAgent()
	ctx.Next()
	end := time.Now().UnixNano() / 1e6
	request := fmt.Sprintf("[path]--> %s [method]--> %s [IP]-->  %s [time]ms-->  %d", p, method, ip, end-start)
	slog.Info(request)
}
