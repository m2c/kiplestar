package utils

import (
	"github.com/kataras/iris/v12"
	"github.com/m2c/kiplestar/commons"
)

func SetXRequestID(ctx iris.Context) iris.Context {
	xRequestID := ctx.Request().Header.Get(commons.X_REQUEST_ID)
	if xRequestID == "" {
		xRequestID = GetUuid()
	}
	ctx.Values().Set(commons.X_REQUEST_ID, xRequestID)
	return ctx
}

func GetXRequestID(ctx iris.Context) string {
	return ctx.Values().GetString(commons.X_REQUEST_ID)
}
