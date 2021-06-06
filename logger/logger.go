package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is the application Logger.
var Logger *zap.SugaredLogger

func init() {
	createLogsDirectory()
	writerSync := getWriterSync()
	encoder := getEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), writerSync, zapcore.DebugLevel),
	)
	zapLogger := zap.New(core).Sugar()
	defer zapLogger.Sync() // flushes buffer, if any
	Logger = zapLogger
	Logger.Debug("Logger initialised")
}

func panicIfNeedbe(e error) {
	if e != nil {
		panic(e)
	}
}

func createLogsDirectory() {
	path, err := os.Getwd()
	panicIfNeedbe(err)
	_, err = os.Stat(fmt.Sprintf("%s/logs", path))
	if os.IsNotExist(err) {
		_ = os.Mkdir("logs", os.ModePerm)
	}
}

func getWriterSync() zapcore.WriteSyncer {
	// path, err := os.Getwd()
	// panicIfNeedbe(err)
	// file, err := os.OpenFile(path+"/logs/nursery.log", os.O_APPEND|os.O_CREATE, 0644)
	// panicIfNeedbe(err)
	// return zapcore.AddSync(file)
	return zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   "logs/nursery.log",
			MaxSize:    1, // megabytes
			MaxBackups: 3,
			MaxAge:     1, // days
			Compress:   true,
		})
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //TODO: UTC time encoder?
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
