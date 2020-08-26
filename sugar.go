package log

import (
	"fmt"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	sugarRootPath = "./"
	sugarLogFile  = "output.log"
	sugarLogLevel = "debug"

	sugar *ZapLoggerWapper
)

func init() {
	sugar = NewZapLoggerWrapper(initSugarLogger().Sugar())
}

func initSugarLogger() *zap.Logger {
	if !path.IsAbs(sugarLogFile) {
		sugarLogFile = path.Join(sugarRootPath, sugarLogFile)
	}
	fmt.Printf("Log path : %s\n", sugarLogFile)
	hook := lumberjack.Logger{
		Filename:   sugarLogFile,
		MaxSize:    500,
		MaxBackups: 100,
		MaxAge:     1000,
		Compress:   true,
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "_time",
		LevelKey:       "_level",
		NameKey:        "_logger",
		CallerKey:      "_caller",
		MessageKey:     "_message",
		StacktraceKey:  "_stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(stringToXZapLoggerLevel(sugarLogLevel))
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(&hook), atomicLevel)
	logger := zap.New(core)
	return logger
}

func stringToXZapLoggerLevel(level string) zapcore.Level {
	lower := strings.ToLower(level)
	switch lower {
	case "info":
		return zap.InfoLevel
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}

// With return a new Logger with context
func With(ctx ...interface{}) Logger {
	return NewZapLoggerWrapper(sugar.logger.With(ctx...))
}

// Debug ...
func Debug(msg string, ctx ...interface{}) {
	sugar.logger.Debugw(msg, ctx...)
}

// Info ...
func Info(msg string, ctx ...interface{}) {
	sugar.logger.Infow(msg, ctx...)
}

// Warn ...
func Warn(msg string, ctx ...interface{}) {
	sugar.logger.Warnw(msg, ctx...)
}

// Error ...
func Error(msg string, ctx ...interface{}) {
	sugar.logger.Errorw(msg, ctx...)
}
