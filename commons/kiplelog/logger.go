package kiplelog

import (
	"github.com/m2c/kiplestar/commons"
	"go.uber.org/zap"
)

type KipleLog struct {
	Logger     *zap.SugaredLogger
	XRequestId string
}

func NewLogger(xid string, logger *zap.SugaredLogger) *KipleLog {
	return &KipleLog{
		Logger:     logger,
		XRequestId: xid,
	}
}

func (log *KipleLog) Prefix() string {
	return commons.X_REQUEST_ID + ":" + log.XRequestId
}

func (log *KipleLog) Print(v ...interface{}) {
	log.Logger.Info(log.Prefix()+"\t", v)
}

func (log *KipleLog) Printf(format string, args ...interface{}) {
	log.Logger.Infof(log.Prefix()+"\t"+format, args)
}

func (log *KipleLog) Info(v ...interface{}) {
	log.Logger.Info(log.Prefix(), v)
}

func (log *KipleLog) Infof(format string, args ...interface{}) {
	log.Logger.Infof(log.Prefix()+"\t"+format, args)
}

func (log *KipleLog) Debug(v ...interface{}) {
	log.Logger.Debug(log.Prefix(), v)
}

func (log *KipleLog) Debugf(format string, args ...interface{}) {
	log.Logger.Debugf(log.Prefix()+"\t"+format, args)
}

func (log *KipleLog) Error(v ...interface{}) {
	log.Logger.Error(log.Prefix(), v)
}

func (log *KipleLog) Errorf(format string, args ...interface{}) {
	log.Logger.Errorf(log.Prefix()+"\t"+format, args)
}
