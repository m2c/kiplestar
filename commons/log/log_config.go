package slog

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/m2c/kiplestar/commons"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

var Log *zap.SugaredLogger
var Slog LogConfig

type LogConfig struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	FileName string `yaml:"filename"`
}

func InitLogger(logConfig LogConfig, app *iris.Application) {
	encoder := getEncoder()
	/*logLevel := zap.DebugLevel*/
	/*level:=logConfig.Level*/
	/*switch level {
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
	}*/
	//info level
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	infoWriter := getLogWriter(logConfig.Path, logConfig.FileName+"-info")
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	//error level
	errorWriter := getLogWriter(logConfig.Path, logConfig.FileName+"-error")
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	// develop mode
	caller := zap.AddCaller()
	// open the code line
	development := zap.Development()
	logger := zap.New(core, caller, development, zap.AddCallerSkip(1))
	Log = logger.Sugar()

	//set iris log level
	if strings.EqualFold(logConfig.Level, "INFO") {
		app.Logger().SetLevel("info")
		app.Logger().SetOutput(infoWriter)
	}
	if strings.EqualFold(logConfig.Level, "Debug") {
		app.Logger().SetLevel("debug")
		app.Logger().SetOutput(os.Stdout)
		p := pprof.New()
		app.Get("/debug/pprof", p)
		app.Get("/debug/pprof/{action:path}", p)
	}
}

/**
 * time format
 */
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05]"))
}

/**
 * get zap log encoder
 */
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	return zapcore.NewConsoleEncoder(encoderConfig)
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
func (slf *LogConfig) Print(v ...interface{}) {
	Log.Info(v)
}
func Info(args ...interface{}) {
	Log.Info(args)
}
func Infof(template string, args ...interface{}) {
	Log.Infof(template, args...)
}
func Debug(args ...interface{}) {
	Log.Debug(args)
}
func Debugf(template string, args ...interface{}) {
	Log.Debugf(template, args...)
}
func Error(args ...interface{}) {
	Log.Error(args)
}
func Errorf(template string, args ...interface{}) {
	Log.Errorf(template, args...)
}

func DebugfCtx(c iris.Context, template string, args ...interface{}) {
	id := c.Values().GetString(commons.X_REQUEST_ID)
	Log.Debugf(commons.X_REQUEST_ID+":"+id+", "+template, args...)
}

func InfofCtx(c iris.Context, template string, args ...interface{}) {
	id := c.Values().GetString(commons.X_REQUEST_ID)
	Log.Infof(commons.X_REQUEST_ID+":"+id+", "+template, args...)
}

func ErrorfCtx(c iris.Context, template string, args ...interface{}) {
	id := c.Values().GetString(commons.X_REQUEST_ID)
	Log.Errorf(commons.X_REQUEST_ID+":"+id+", "+template, args...)
}
