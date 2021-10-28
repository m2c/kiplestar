package slog

import (
	"go.uber.org/zap"
)

type KipleLogger struct {
	Logger     *zap.SugaredLogger
	XRequestId string
}

func NewLogger(xid string, logger *zap.SugaredLogger) *KipleLogger {
	return &KipleLogger{
		Logger:     logger,
		XRequestId: xid,
	}
}

// for db logger
func (log *KipleLogger) Print(v ...interface{}) {
	log.Logger.Info(v...)
}

func (log *KipleLogger) Info(v ...interface{}) {
	log.Logger.Info(v...)
}

func (log *KipleLogger) Infof(format string, args ...interface{}) {
	log.Logger.Infof(format, args...)
}

func (log *KipleLogger) Debug(v ...interface{}) {
	log.Logger.Debug(v...)
}

func (log *KipleLogger) Debugf(format string, args ...interface{}) {
	log.Logger.Debugf(format, args...)
}

func (log *KipleLogger) Error(v ...interface{}) {
	log.Logger.Error(v...)
}

func (log *KipleLogger) Errorf(format string, args ...interface{}) {
	log.Logger.Errorf(format, args...)
}

func (log *KipleLogger) Warn(v ...interface{}) {
	log.Logger.Warn(v...)
}

func (log *KipleLogger) Warnf(format string, args ...interface{}) {
	log.Logger.Warnf(format, args...)
}
