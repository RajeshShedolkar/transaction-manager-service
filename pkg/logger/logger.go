package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(env string) error {
	var err error

	if env == "prod" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = newDevLogger()
	}

	if err != nil {
		return err
	}
	zap.ReplaceGlobals(Log)
	return nil
}

func newDevLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg.Build()
}
