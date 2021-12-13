package middleware

import (
	"bytes"
	"fmt"
	"github.com/kataras/iris/v12"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/commons/utils"
	"io/ioutil"
	"runtime"
	"strings"
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

	// iris.WithoutBodyConsumptionOnUnmarshal is removed out of kiplestar, so read body by hand here.
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		slog.Infof("ReadAll body failed: %s", err.Error())
	} else {
		ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
		bodyLen := len(body)
		if bodyLen > 0 && bodyLen < 500 {
			// format body to one line for align log system
			slog.Infof("log http request body: %s", strings.Replace(utils.SensitiveFilter(string(body)), "\n", " ", -1))
		}
	}

	// calculate cost time
	ctx.Next()
}
