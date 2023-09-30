package logger

import (
	"log"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(log_path string) *zap.Logger {

	if len(log_path) == 0 {
		log.Fatal("Empty path for local log file")
	}

	writerSyncer := getLogWriter(log_path)

	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())

	defer logger.Sync()

	return logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(log_path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   log_path,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
