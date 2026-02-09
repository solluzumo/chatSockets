package app

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger() *zap.Logger {
	serviceName := os.Getenv("SERVICE_NAME")
	logsDirName := fmt.Sprintf("logs/%s", serviceName)

	cfg := zap.NewProductionConfig()
	os.MkdirAll(logsDirName, os.ModePerm)
	cfg.OutputPaths = []string{"stdout", logsDirName + "/service.log"}
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//СЭМПЛИРОВАНИЕ ЛОГОВ
	cfg.Sampling = &zap.SamplingConfig{
		Initial:    50,  // первые 50 сообщений за 1 сек — все логируются
		Thereafter: 100, // после этого — только каждое 100-е
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("не получилось создать логгер: %v", err)
	}

	return logger
}
