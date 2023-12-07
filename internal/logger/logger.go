package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var zapLevels = map[string]zapcore.Level{
	"warn":  zapcore.WarnLevel,
	"info":  zapcore.InfoLevel,
	"debug": zapcore.DebugLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
	"panic": zapcore.PanicLevel,
}

func CreateLogger() *zap.Logger {
	var encoderCfg zapcore.EncoderConfig
	if os.Getenv("APP_MODE") == "production" {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.MessageKey = "message"
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "caller"
	encoderCfg.NameKey = "name"

	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// Create writers for file and console
	consoleWriter := zapcore.AddSync(os.Stdout)

	l := &lumberjack.Logger{
		Filename: fmt.Sprintf(
			"%saplication-%s.log",
			os.Getenv("LOG_FILE_PATH"),
			time.Now().Format("2006-01-02T15-04-05"),
		), // Or any other path
		MaxSize:    500,  // MB; after this size, a new log file is created
		MaxBackups: 3,    // Number of backups to keep
		MaxAge:     28,   // Days
		Compress:   true, // Compress the backups using gzip
	}

	fileWriter := zapcore.AddSync(l)

	// Create cores for writing to the file and console
	fileCore := zapcore.NewCore(
		fileEncoder,
		fileWriter,
		zapLevels[os.Getenv("LOG_LEVEL")],
	)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		consoleWriter,
		zapLevels[os.Getenv("LOG_LEVEL")],
	)

	// Combine cores
	core := zapcore.NewTee(fileCore, consoleCore)

	// Create the logger with additional context information (caller, stack trace)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
