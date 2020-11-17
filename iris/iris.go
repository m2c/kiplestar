package iris

import (
	"errors"
	"fmt"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/m2c/kiplestar/commons"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/config"
	"github.com/m2c/kiplestar/middleware"
	"net/http"
)

type App struct {
	app *iris.Application
}

func (slf *App) Default() {
	slf.app = iris.New()
	//register middleware
	slf.app.UseGlobal(middleware.GlobalRecover, middleware.LoggerHandler)
	//global error handling
	slf.app.OnAnyErrorCode(func(ctx iris.Context) {
		_, _ = ctx.JSON(commons.BuildFailedWithMsg(commons.ResponseCode(ctx.GetStatusCode()), http.StatusText(ctx.GetStatusCode())))
	})
	slf.initServerLog()
}

func (slf *App) New() {
	slf.app = iris.New()
	//global error handling
	slf.app.OnAnyErrorCode(func(ctx iris.Context) {
		_, _ = ctx.JSON(commons.BuildFailedWithMsg(commons.UnKnowError, ctx.Values().GetString("message")))
	})
	slf.initServerLog()
}

//init server log
func (slf *App) initServerLog() {
	slog.Slog = slog.LogConfig{Level: config.SC.SConfigure.LogLevel, Path: config.SC.SConfigure.LogPath, FileName: config.SC.SConfigure.LogName}
	slog.InitLogger(slog.Slog, slf.app)

}

//set middleware
func (slf *App) SetGlobalMiddleware(handlers ...context.Handler) {
	slf.app.UseGlobal(handlers...)
}

//set middleware
func (slf *App) SetMiddleware(handlers ...context.Handler) {
	slf.app.Use(handlers...)
}

//get Iris App
func (slf *App) GetIrisApp() *iris.Application {
	return slf.app
}

func (slf *App) Party(relativePath string, handlers ...context.Handler) {
	slf.app.Party(relativePath, handlers...)
}
func (slf *App) Post(relativePath string, handlers ...context.Handler) {
	slf.app.Post(relativePath, handlers...)
}
func (slf *App) Get(relativePath string, handlers ...context.Handler) {
	slf.app.Get(relativePath, handlers...)
}

//start server,
func (slf *App) Start(params ...iris.Configurator) error {
	server := fmt.Sprintf("%s:%d", config.SC.SConfigure.Addr, config.SC.SConfigure.Port)
	if slf.app == nil {
		return errors.New("Server not init")
	}
	//go slf.app.Run(iris.Addr(server))
	swaggerConfig := &swagger.Config{
		URL: fmt.Sprintf("./swagger/doc.json"), //The url pointing to API definition
	}
	slf.app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(swaggerConfig, swaggerFiles.Handler))
	p := pprof.New()
	slf.app.Get("/debug/pprof", p)
	slf.app.Get("/debug/pprof/{action:path}", p)
	params = append(params, iris.WithoutStartupLog)
	if config.SC.SConfigure.Profile != "prod" {
		params = append(params, iris.WithoutBodyConsumptionOnUnmarshal)
	}
	return slf.app.Run(iris.Addr(server), params...)
}
