package internal

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Loggger = func() *zap.Logger {
	var (
		logger *zap.Logger
		err    error
	)
	if os.Getenv("TZF_DEBUG") == "true" {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		logger, err = cfg.Build()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	return logger
}()
