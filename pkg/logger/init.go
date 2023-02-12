package logger

import (
	"fmt"

	"github.com/Nikolay-Yakushev/lychee/config"
	"go.uber.org/zap"
)



func New(cfg *config.Config) (*zap.Logger, error) {
	zapCfg := zap.NewProductionConfig()

	lvl, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("Log level parse error. Reason: %w", err)
	}

	zapCfg.Level = lvl
	logger, err := zapCfg.Build()
	if err != nil {
		err := fmt.Errorf("Failed zap logger init: %w", err)
		return nil, err
	}
	return logger, nil
}
	
