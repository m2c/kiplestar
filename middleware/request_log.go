package middleware

import (
	"github.com/kataras/iris/v12"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/commons/utils"
	"github.com/m2c/kiplestar/config"
	"strings"
	"time"
)

func LoggerHandler(ctx iris.Context) {
	ctx = utils.SetXRequestID(ctx)
	p := ctx.Request().URL.Path
	method := ctx.Request().Method
	start := time.Now().UnixNano() / 1e6
	ip := ctx.Request().RemoteAddr

	ctx.Next()
	end := time.Now().UnixNano() / 1e6
	slog.InfofCtx(ctx, "[path]--> %s [method]--> %s [IP]-->  %s [time]ms-->  %d", p, method, ip, end-start)
	if config.SC.SConfigure.Profile != "prod" {
		body, err := ctx.GetBody()
		if err != nil {
			return
		}
		// format body to one line for aliyun log system
		slog.InfofCtx(ctx, "debug log http request body %s", strings.Replace(string(body), "\n", " ", -1))
	}
}
