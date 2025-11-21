package config

import (
	"elderly-care-backend/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() {

	var logger *zap.Logger

	switch Config.Server.Profile {
	case "dev":
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			//EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
			EncodeLevel:    zapcore.CapitalColorLevelEncoder, //这里可以指定颜色
			EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		}
		atom := zap.NewAtomicLevelAt(zap.DebugLevel)
		config := zap.Config{
			Level:            atom,
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		logger, _ = config.Build()
	case "prod":
		writer := zapcore.AddSync(logFile)
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core := zapcore.NewCore(encoder, writer, zap.ErrorLevel)
		logger = zap.New(core)
	}
	global.Logger = logger
}
