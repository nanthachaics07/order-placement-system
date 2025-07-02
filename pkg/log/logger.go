package log

import (
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ? Logger นี้มาจาก universal-lib ทีที่ผมใช้ในโปรเจคอื่นๆ
type Logger struct {
	zap   *zap.Logger
	sugar *zap.SugaredLogger
}

type Field struct {
	key string
	val interface{}
}

var (
	instance *Logger
	once     sync.Once
)

func Init(env string) {
	if env != "dev" && env != "prod" {
		fmt.Println("Invalid environment specified, defaulting to 'dev' but please use 'dev' or 'prod'")
		env = "dev"
	}

	once.Do(func() {
		var cfg zap.Config
		switch strings.ToLower(env) {
		case "dev":
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
			cfg.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
			cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
			cfg.EncoderConfig.ConsoleSeparator = " | "
		case "prod":
			cfg = zap.Config{
				Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
				Development:      false,
				Encoding:         "json",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
				EncoderConfig: zapcore.EncoderConfig{
					TimeKey:    "time",
					LevelKey:   "level",
					NameKey:    "logger",
					CallerKey:  "caller",
					MessageKey: "message",
					// StacktraceKey:  "stacktrace",
					LineEnding:     zapcore.DefaultLineEnding,
					EncodeLevel:    zapcore.LowercaseLevelEncoder,
					EncodeTime:     zapcore.ISO8601TimeEncoder,
					EncodeDuration: zapcore.SecondsDurationEncoder,
					EncodeCaller:   zapcore.FullCallerEncoder,
				},
			}
		default:
			cfg = zap.NewDevelopmentConfig()
		}

		z, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}

		instance = &Logger{
			zap:   z,
			sugar: z.Sugar(),
		}
	})
}

func Get() *Logger {
	if instance == nil {
		panic("logger not initialized, call logger.Init(env) first")
	}
	return instance
}

func Info(msg string) {
	Get().zap.Info(msg)
}

func Debug(msg string) {
	Get().zap.Debug(msg)
}

func Error(msg string) {
	Get().zap.Error(msg)
}

func Warn(msg string) {
	Get().zap.Warn(msg)
}

func Fatal(msg string) {
	Get().zap.Fatal(msg)
}

func Infof(msg string, args ...interface{}) {
	logWithFields("info", msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	logWithFields("debug", msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	logWithFields("error", msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	logWithFields("warn", msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	logWithFields("fatal", msg, args...)
}

func logWithFields(level string, msg string, args ...interface{}) {
	fields := make([]zap.Field, 0)
	others := make([]interface{}, 0)

	for _, arg := range args {
		switch v := arg.(type) {
		case Field:
			fields = append(fields, zap.Any(v.key, v.val))
		case error:
			fields = append(fields, zap.Error(v))
		case string:
			fields = append(fields, zap.String("message", v))
		case int:
			fields = append(fields, zap.Int("value", v))
		case bool:
			fields = append(fields, zap.Bool("flag", v))
		case float64:
			fields = append(fields, zap.Float64("value", v))
		case float32:
			fields = append(fields, zap.Float32("value", v))
		case []byte:
			fields = append(fields, zap.ByteString("data", v))
		case map[string]interface{}:
			for k, val := range v {
				fields = append(fields, zap.Any(k, val))
			}
		case []interface{}:
			for _, item := range v {
				fields = append(fields, zap.Any("item", item))
			}
		case nil:
			fields = append(fields, zap.String("nil", "nil"))

		default:
			others = append(others, v)
		}
	}

	if len(others) > 0 {
		msg = msg + " | " + fmt.Sprint(others...)
	}

	logger := Get().zap.WithOptions(zap.AddCallerSkip(1))

	switch level {
	case "info":
		logger.Info(msg, fields...)
	case "debug":
		logger.Debug(msg, fields...)
	case "error":
		logger.Error(msg, fields...)
	case "warn":
		logger.Warn(msg, fields...)
	case "fatal":
		logger.Fatal(msg, fields...)
	}
}

func E(err error) Field {
	return Field{key: "error", val: err}
}

func S(key, val string) Field {
	return Field{key: key, val: val}
}

func AtoS(key string, val any) Field {
	return Field{key: key, val: fmt.Sprintf("%v", val)}
}

func Any(key string, val interface{}) Field {
	return Field{key: key, val: val}
}

func Sync() error {
	if instance != nil {
		return instance.zap.Sync()
	}
	return nil
}
