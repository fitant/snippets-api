package utils

import (
	"github.com/fitant/xbin-api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


var logLevels map[string]zapcore.Level = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info": zap.InfoLevel,
	"warn": zap.WarnLevel,
	"error": zap.ErrorLevel,
}

func InitLogger(cfg *config.Config) *zap.Logger {
	lgrCfg := zap.NewProductionConfig()
	if cfg.App.GetEnv() == config.Environments["dev"] {
		lgrCfg.Development = true
	}
	lgrCfg.Level = zap.NewAtomicLevelAt(logLevels[cfg.App.GetLogLevel()])

	lgr, err := lgrCfg.Build()
	if err != nil {
		panic(err)
	}
	return lgr
}
