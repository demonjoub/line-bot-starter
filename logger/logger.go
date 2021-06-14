package logger

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Logger type
type Logger struct {
	*zap.Logger
}

// RequestInfo .
type RequestInfo struct {
	Action        string
	TraceID       string
	ParentID      string
	SpanID        string
	RemoteAddress string
	Tag           string
	Msg           string
}

// Config defines logger instance config
type Config struct {
	Environment string
	Level       string
}

var logger *zap.Logger = zap.NewExample()

func InitLogger(cf *Config) error {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	var cfg zap.Config
	if cf.Environment == "production" {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout"}
		if cf.Level == "debug" {
			cfg.Level.SetLevel(zapcore.DebugLevel)
		}
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.EncoderConfig = ec

	l, err := cfg.Build()
	if err != nil {
		return errors.WithStack(err)
	}

	hn, _ := os.Hostname()
	if hn == "" {
		hn = "unknown"
	}

	logger = l.With(zap.String("hostname", hn))
	return nil
}

func NewLoggerWithRequestInfo(ri *RequestInfo) Logger {
	return Logger{
		logger.With(
			zap.String("trace_id", ri.TraceID),
			zap.String("parent_id", ri.ParentID),
			zap.String("span_id", ri.SpanID),
			zap.String("remote_address", ri.RemoteAddress),
			zap.String("node_id", ""),
		),
	}
}

func NewRequestLogger(reqID string, apiPath string) Logger {
	return Logger{
		logger.With(
			zap.String("request_id", reqID),
			zap.String("api_path", apiPath),
		),
	}
}

func Info(msg string, fields ...zapcore.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}
