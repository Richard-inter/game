package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func InitLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	Sugar = Logger.Sugar()
}

func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

func GetSugar() *zap.SugaredLogger {
	if Sugar == nil {
		InitLogger()
	}
	return Sugar
}

func SetLevel(level string) {
	var lvl zapcore.Level
	switch level {
	case "debug":
		lvl = zap.DebugLevel
	case "info":
		lvl = zap.InfoLevel
	case "warn":
		lvl = zap.WarnLevel
	case "error":
		lvl = zap.ErrorLevel
	case "fatal":
		lvl = zap.FatalLevel
	default:
		lvl = zap.InfoLevel
	}

	if Logger != nil {
		Logger.Core().Sync()
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level = zap.NewAtomicLevelAt(lvl)

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	Sugar = Logger.Sugar()
}

func SetFormatter(format string) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch format {
	case "json":
		config.Encoding = "json"
	case "text":
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.Encoding = "json"
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	Sugar = Logger.Sugar()
}
