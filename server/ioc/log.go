package ioc

import (
	"go.uber.org/zap"
	"ink-im-server/pkg/logger"
)

func InitLogger() logger.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
