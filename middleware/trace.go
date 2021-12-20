package middleware

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/m2c/kiplestar/commons"
	"github.com/m2c/kiplestar/commons/utils"
	slog "github.com/m2c/kiplestar/commons/log"
	uuid "github.com/satori/go.uuid"
	"time"
	"strconv"
)

//TraceLogger used for record request id
func TraceLogger(ctx iris.Context) {
	requestID := ctx.Request().Header.Get(commons.X_REQUEST_ID)
	if len(requestID) == 0 {
		requestID = uuid.NewV4().String()
	}
	parentSpanIdStr := ctx.Request().Header.Get(commons.X_SPAN_ID)
	if len(parentSpanIdStr) == 0 {
		parentSpanIdStr = "1"
	}
	parentSpanId := int32(utils.StringToInt(parentSpanIdStr,1))
	span := slog.Span{
		ParentSpanID: parentSpanId,
		SubSpanID : 0,
		NextSpanID : parentSpanId,
	}
	traceContext := context.WithValue(ctx.Request().Context(), commons.X_REQUEST_ID, requestID)
	spanContext := context.WithValue(traceContext, commons.X_SPAN_ID, &span)
	newRequest := ctx.Request().WithContext(spanContext)
	ctx.ResetRequest(newRequest)
	path := ctx.Request().URL.Path
	method := ctx.Request().Method
	ip := ctx.Request().RemoteAddr

	slog.InfofStdCtx(spanContext, "rid:%s path:%s method:%s ip:%s start \n", requestID, path, method, ip)
	start := time.Now().UnixNano() / 1e6
	ctx.Next()
	end := time.Now().UnixNano() / 1e6
	slog.Log.Infow(
		"",
		commons.X_REQUEST_ID, requestID,
		commons.X_SPAN_ID, parentSpanId,
		commons.LOG_FIELD_SERVICE, "selank-merchant-service", //define in config file
		commons.LOG_FIELD_RSP_TIME, strconv.FormatInt(end-start, 10),
		commons.LOG_FIELD_URL, ctx.Request().URL,
		commons.LOG_FIELD_METHOD, ctx.Request().Method,
		commons.LOG_FIELD_HEADER, ctx.Request().Header,
		commons.LOG_FIELD_ARGS, "",
		commons.LOG_FIELD_CLIENT, ip,
		commons.LOG_FIELD_HTTP_CODE, ctx.GetStatusCode(),
	)
	//slog.InfofStdCtx(traceContext, "done")
}
