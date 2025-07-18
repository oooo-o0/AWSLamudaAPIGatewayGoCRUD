package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	var cfg zap.Config
	if os.Getenv("ENV") == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}

	Logger = logger.Sugar()
}
