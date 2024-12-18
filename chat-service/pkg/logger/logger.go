package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(env string) {
    var config zap.Config
    if env == "production" {
        config = zap.NewProductionConfig()
    } else {
        config = zap.NewDevelopmentConfig()
    }
    
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    
    var err error
    log, err = config.Build()
    if err != nil {
        panic(err)
    }
}

func Info(msg string, fields ...zap.Field) {
    log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
    log.Fatal(msg, fields...)
}
