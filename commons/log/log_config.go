package slog

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/m2c/kiplestar/commons"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var ZapLogger *zap.Logger
var Log *zap.SugaredLogger
var Slog LogConfig

type LogConfig struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	FileName string `yaml:"filename"`
}

func init() {
	InitLogger(LogConfig{}, nil)
}

func Logger(xid string) *KipleLogger {
	logger := ZapLogger.WithOptions(zap.Fields(
		zap.String(commons.X_REQUEST_ID, xid),
	), zap.AddCallerSkip(0))
	return NewLogger(xid, logger.Sugar())
}

func InitLogger(logConfig LogConfig, app *iris.Application) {
	encoder := getEncoder()
	var writer io.Writer
	if logConfig.FileName != "" {
		writer = io.MultiWriter(os.Stdout, getLogWriter(logConfig.Path, logConfig.FileName))
	} else {
		writer = os.Stdout
	}
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), getLogLevel(logConfig.Level)),
	)
	// develop mode
	caller := zap.AddCaller()
	// open the code line
	development := zap.Development()
	ZapLogger = zap.New(core, caller, development, zap.AddCallerSkip(1))
	Log = ZapLogger.Sugar()

	//set iris log level
	if app != nil {
		app.Logger().SetLevel(logConfig.Level)
		app.Logger().SetOutput(writer)
	}
}

func getLogLevel(level string) zapcore.Level {
	var logLevel zapcore.Level
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}
	return logLevel
}

/**
 * time format
 */
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

/**
 * get zap log encoder
 */
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//encoderConfig.LineEnding = zapcore.DefaultLineEnding
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(logPath, level string) io.Writer {
	logFullPath := path.Join(logPath, level)
	hook, err := rotatelogs.New(
		logFullPath+"-%Y%m%d%H"+".txt",
		// log file split
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func getLogPrefix() string {
	return fmt.Sprintf("%s: %s  ", commons.X_REQUEST_ID, GetLogID())
}

func (slf *LogConfig) Print(v ...interface{}) {
	// Temporarily remove the log xid of db
	Log.Info(v)
}

func Info(args ...interface{}) {
	Log.Info(getLogPrefix(), args)
}

func Infof(template string, args ...interface{}) {
	Log.Infof(getLogPrefix()+template, args...)
}

func Debug(args ...interface{}) {
	Log.Debug(getLogPrefix(), args)
}

func Debugf(template string, args ...interface{}) {
	Log.Debugf(getLogPrefix()+template, args...)
}

func Error(args ...interface{}) {
	Log.Error(getLogPrefix(), args)
}

func Errorf(template string, args ...interface{}) {
	if runtime.GOOS == "linux" {
		msg := strings.ReplaceAll(fmt.Sprintf(getLogPrefix()+template, args...), "\n", "\\d")
		Log.Error(msg)
		return
	}
	Log.Errorf(getLogPrefix()+template, args...)
}
func DebugfStdCtx(c context.Context, template string, args ...interface{}) {
	xid := c.Value(commons.X_REQUEST_ID).(string)
	Log.Debugf(commons.X_REQUEST_ID+": "+xid+"  "+template, args...)
}

func InfofStdCtx(c context.Context, template string, args ...interface{}) {
	xid := c.Value(commons.X_REQUEST_ID).(string)
	Log.Infof(commons.X_REQUEST_ID+": "+xid+"  "+template, args...)
}

func ErrorfStdCtx(c context.Context, template string, args ...interface{}) {
	xid := c.Value(commons.X_REQUEST_ID).(string)
	Log.Errorf(commons.X_REQUEST_ID+": "+xid+"  "+template, args...)
}

func DebugfCtx(c iris.Context, template string, args ...interface{}) {
	xid := c.Values().GetString(commons.X_REQUEST_ID)
	Log.Debugf(commons.X_REQUEST_ID+": "+xid+"  "+template, args...)
}

func InfofCtx(c iris.Context, template string, args ...interface{}) {
	xid := c.Values().GetString(commons.X_REQUEST_ID)
	Log.Infof(commons.X_REQUEST_ID+": "+xid+"  "+template, args...)
}

func ErrorfCtx(c iris.Context, template string, args ...interface{}) {
	xid := c.Values().GetString(commons.X_REQUEST_ID)
	Log.Errorf(commons.X_REQUEST_ID+": "+xid+"  "+template, args...)
}

func InfowCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	Log.Infow(msg, append(getTrace(c), keysAndValues...)...)
}

func ErrorwCtx(c context.Context, msg string, keysAndValues ...interface{}) {
	Log.Errorw(msg, append(getTrace(c), keysAndValues...)...)
}

func getTrace(c context.Context) []interface{} {
	span := (c.Value(commons.X_SPAN_ID)).(*Span)
	return []interface{}{
		commons.X_REQUEST_ID, c.Value(commons.X_REQUEST_ID).(string),
		commons.X_SPAN_ID, GetSubSpanID(c, span),
	}
}
