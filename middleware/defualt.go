package middleware

import (
	"fmt"
	"github.com/kataras/iris/v12"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/commons/utils"
	"github.com/m2c/kiplestar/config"
	"runtime"
	"strings"
	"time"
)

func Default(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}

			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}

				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}

			// when stack finishes
			logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
			logMessage += fmt.Sprintf("Trace: %s", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			//ctx.Application().Logger().Error(logMessage)
			slog.Error(logMessage)
			ctx.StatusCode(500)
			ctx.StopExecution()
		}
	}()

	ctx = utils.SetXRequestID(ctx)
	p := ctx.Request().URL.Path
	method := ctx.Request().Method
	start := time.Now().UnixNano() / 1e6
	ip := ctx.Request().RemoteAddr
	slog.SetLogID(utils.GetXRequestID(ctx))

	ctx.Next()
	end := time.Now().UnixNano() / 1e6
	slog.Infof("[path]--> %s [method]--> %s [IP]-->  %s [time]ms-->  %d", p, method, ip, end-start)
	if config.SC.SConfigure.Profile != "prod" {
		body, err := ctx.GetBody()
		if err != nil {
			return
		}
		if len(body) > 0 {
			// format body to one line for aliyun log system
			slog.Infof("log http request body: %s", strings.Replace(utils.SensitiveFilter(string(body)), "\n", " ", -1))
		}
	}
}
