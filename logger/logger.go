package logger

import (
	"core/domain"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {

	var err error
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	config.EncoderConfig = encoderConfig
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) *domain.Log {
	log.Info(message, fields...)
	a := log.Check(zapcore.InfoLevel, message)
	lg := domain.Log{
		TimeStamp: a.Time,
		Info:      a.Message,
		Function:  a.Caller.Function,
	}
	return &lg
}

func Debug(message string, fields ...zap.Field) *domain.Log {
	log.Debug(message, fields...)
	a := log.Check(zapcore.InfoLevel, message)
	lg := domain.Log{
		TimeStamp: a.Time,
		Info:      a.Message,
		Function:  a.Caller.Function,
	}
	return &lg
}

func Error(message string, fields ...zap.Field) *domain.Log {
	log.Error(message, fields...)
	a := log.Check(zapcore.InfoLevel, message)
	lg := domain.Log{
		TimeStamp: a.Time,
		Info:      a.Message,
		Function:  a.Caller.Function,
	}
	return &lg
}
