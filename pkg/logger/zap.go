package logger

import (
	"context"
	"log"
	"micro/adapter/elk"
	"micro/config"
	"micro/pkg/oauth"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.elastic.co/ecszap"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func InitGlobalLogger(lc fx.Lifecycle, logstash *elk.LogStash) error {
	logger := configLogger(logstash)
	zap.ReplaceGlobals(logger)

	lc.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			if err := zap.L().Sync(); err != nil {
				log.Println("logger failed to sync:", err)
			}
			return nil
		},
	})
	return nil
}

func configLogger(el *elk.LogStash) *zap.Logger {
	logLevel := getLogLevel()

	elkZapCore := ecszap.NewCore(
		ecszap.NewDefaultEncoderConfig(),
		zapcore.AddSync(el),
		logLevel,
	)

	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	terminalZapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		logLevel,
	)

	core := zapcore.NewTee(elkZapCore, terminalZapCore)
	logger := zap.New(core, zap.AddCaller())
	return logger.With(zap.String("service", config.C().Service.Name))
}

func getLogLevel() zapcore.Level {
	if config.C().Debug {
		return zap.DebugLevel
	}
	return zap.InfoLevel
}

type zapCtxKeyType string

const (
	zapCtxKey zapCtxKeyType = "zap"
)

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(zapCtxKey).(*zap.Logger)
	if !ok {
		panic("cannot extract logger from context")
	}
	return logger
}

func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, zapCtxKey, logger)
}

func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		z := zap.L()
		if userID, err := oauth.UserIDFromContext(ctx); err == nil {
			z = z.With(zap.String("user.id", userID))
		}
		if userPhone, err := oauth.UserMobileFromContext(ctx); err == nil {
			z = z.With(zap.String("user.mobile", userPhone))
		}
		traceID := opentracing.SpanFromContext(ctx).Context().(jaeger.SpanContext).TraceID().String()
		t0 := time.Now()
		z = z.With(
			zap.String("trace.id", traceID),
			zap.String("method", info.FullMethod))
		resp, err := handler(ToContext(ctx, z), req)
		t1 := time.Now()
		z.Info(
			info.FullMethod,
			zap.Time("trace.start_time", t0),
			zap.Time("trace.end_time", t1),
			zap.Duration("trace.duration", t1.Sub(t0)))
		return resp, err
	}
}
