package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap() (*zap.SugaredLogger, func() error, error) {
	if err := os.MkdirAll("_out/logs", 0755); err != nil {
		return nil, nil, fmt.Errorf("mkdir logs: %s", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	filePath := filepath.Join("_out/logs", fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %s", err)
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15-04-05.000000")

	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.DebugLevel),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.WarnLevel),
	)

	return logger.Sugar(), logFile.Close, nil
}
