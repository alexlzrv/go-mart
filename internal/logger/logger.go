package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func LogInitializer(level string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("attempt to parse logger level failed - %v", err)
	}

	cfg.Level = lvl

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("attempt to initialized logger failed - %v", err)
	}

	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)

	return logger.Sugar(), nil
}
