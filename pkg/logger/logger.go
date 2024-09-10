package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"doitsolutions.vn/pkg/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry
	Core() zapcore.Core
	DPanic(msg string, fields ...zapcore.Field)
	Debug(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Level() zapcore.Level
	Log(lvl zapcore.Level, msg string, fields ...zapcore.Field)
	Name() string
	Named(s string) *zap.Logger
	Panic(msg string, fields ...zapcore.Field)
	Sugar() *zap.SugaredLogger
	Sync() error
	Warn(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) *zap.Logger
	WithLazy(fields ...zapcore.Field) *zap.Logger
	WithOptions(opts ...zap.Option) *zap.Logger
}
type logger struct {
	*zap.Logger
}

func getZapLogLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
func getStringLogLevel(lvl zapcore.Level) string {
	switch lvl {
	case zap.DebugLevel:
		return "debug"
	case zap.InfoLevel:
		return "info"
	case zap.WarnLevel:
		return "warn"
	case zap.ErrorLevel:
		return "error"
	default:
		return "info"
	}
}

var LogLevels = [4]zapcore.Level{
	zap.DebugLevel,
	zap.InfoLevel,
	zap.WarnLevel,
	zap.ErrorLevel,
}

func New(cfg configs.LogConfig) Logger {
	logLevel := getZapLogLevel(cfg.LogLevel)
	lumberWriter := &lumberjack.Logger{
		Filename:   getLogFile(cfg.Dir, getStringLogLevel(logLevel)),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	cores := make([]zapcore.Core, 0)
	for _, v := range LogLevels {
		lb := lumberjack.Logger{
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		if logLevel <= v {
			cores = append(cores, createCore(&lb, cfg.Dir, v))
		}
	}
	w := zapcore.AddSync(lumberWriter)
	w = zapcore.NewMultiWriteSyncer(w, zapcore.AddSync(os.Stdout))
	encoder := getEncoderLog()
	cores = append(cores, zapcore.NewCore(encoder, w, logLevel))
	core := zapcore.NewTee(cores...)
	return &logger{
		zap.New(core),
	}
}

func createCore(lumberWriter *lumberjack.Logger, dir string, level zapcore.Level) zapcore.Core {
	lumberWriter.Filename = getLogFile(dir, getStringLogLevel(level))
	w := zapcore.AddSync(lumberWriter)
	encoder := getEncoderLog()
	return zapcore.NewCore(encoder, w, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == level
	}))
}

func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encodeConfig)
}
func getLogFile(dir string, lvl string) string {
	currentTime := time.Now()
	filename := fmt.Sprintf("%s.%s.log", currentTime.Format("2006-01-02"), lvl)
	return filepath.Join(dir, filename)
}
