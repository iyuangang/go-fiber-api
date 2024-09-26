package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
	once sync.Once
)

func InitLogger() {
	once.Do(func() {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder

		fileEncoder := zapcore.NewJSONEncoder(config)
		consoleEncoder := zapcore.NewConsoleEncoder(config)

		logFile, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer := zapcore.AddSync(logFile)

		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, writer, zapcore.InfoLevel),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		)

		Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
}

func GetLogger() *zap.Logger {
	if Log == nil {
		InitLogger()
	}
	return Log
}
