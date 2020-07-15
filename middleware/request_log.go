package middleware

import (
	"fmt"
	"github.com/kataras/iris/v12"
	slog "kiple_star/commons/log"
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
	time := end - start
	request := fmt.Sprintf("[path]--> %s [method]--> %s [IP]-->  %s [time]ms-->  %d", p, method, ip, time)
	slog.Info(request)
	slog.Error(request)
}
