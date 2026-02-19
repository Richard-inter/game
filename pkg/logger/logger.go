package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	TimeKey    = "timestamp"
	JSONFormat = "json"
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func InitLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = TimeKey
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
	case DebugLevel:
		lvl = zap.DebugLevel
	case InfoLevel:
		lvl = zap.InfoLevel
	case WarnLevel:
		lvl = zap.WarnLevel
	case ErrorLevel:
		lvl = zap.ErrorLevel
	case FatalLevel:
		lvl = zap.FatalLevel
	default:
		lvl = zap.InfoLevel
	}

	if Logger != nil {
		_ = Logger.Sync()
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = TimeKey
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
	config.EncoderConfig.TimeKey = TimeKey
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch format {
	case JSONFormat:
		config.Encoding = JSONFormat
	case "text":
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.Encoding = JSONFormat
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	Sugar = Logger.Sugar()
}
