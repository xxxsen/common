package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logInst, _ = zap.NewProduction()

func stringToLevel(lv string) zapcore.Level {
	lolv := strings.ToLower(lv)
	switch lolv {
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	case "warn":
		return zapcore.WarnLevel
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		panic("unsupport level:" + lv)
	}
}

func Init(file string, lv string, maxRotate int, maxSize int, maxKeepDays int, withConsole bool) *zap.Logger {
	if len(lv) == 0 || len(file) == 0 {
		lv = "debug"
		withConsole = true
	}
	levelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.CapitalString())
	}
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	callerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}
	nameEncoder := func(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(loggerName)
	}
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		CallerKey:        "caller",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeTime:       timeEncoder,
		EncodeCaller:     callerEncoder,
		EncodeName:       nameEncoder,
		EncodeLevel:      levelEncoder,
		EncodeDuration:   zapcore.MillisDurationEncoder,
		ConsoleSeparator: "|",
	})
	synclist := make([]zapcore.WriteSyncer, 0, 2)
	if len(file) != 0 && maxSize > 0 {
		sz := maxSize / 1024 / 1024
		if sz == 0 {
			sz = 1
		}
		logger := &lumberjack.Logger{
			// 日志输出文件路径
			Filename:   file,
			MaxSize:    sz, // megabytes
			MaxBackups: maxRotate,
			MaxAge:     maxKeepDays, //days
			Compress:   false,       // disabled by default
		}
		synclist = append(synclist, zapcore.AddSync(logger))
	}
	if withConsole {
		synclist = append(synclist, zapcore.AddSync(os.Stderr))
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(synclist...),
		stringToLevel(lv),
	)
	logInst = zap.New(core, zap.WithCaller(true))
	return logInst
}

func Logger() *zap.Logger {
	return logInst
}
