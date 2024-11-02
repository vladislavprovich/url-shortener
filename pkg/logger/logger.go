package logger

import "go.uber.org/zap"

func NewLogger(level string) *zap.Logger {
	var cfg zap.Config
	// TODO fix vars
	if level == "prod" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}

	return logger
}
