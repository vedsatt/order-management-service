package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type loggerKey string

const (
	key loggerKey = "logger"
)

type Logger struct {
	z *zap.Logger
}

func New(ctx context.Context, env string) (context.Context, error) {
	loggerCfg := zap.NewProductionConfig()

	loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if env == "dev" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = logger.Sync(); err != nil {
			logger.Error("failed to syncronize logger: %w", zap.Error(err))
		}
	}()

	ctx = context.WithValue(ctx, key, &Logger{logger})

	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.z.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.z.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.z.Error(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.z.Fatal(msg, fields...)
}
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.z.Debug(msg, fields...)
}

func LoggerInterceptor(rootCtx context.Context) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		logger := GetLoggerFromCtx(rootCtx)
		ctx = context.WithValue(ctx, key, GetLoggerFromCtx(rootCtx))

		logger.Info(ctx,
			"incoming request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.Error(ctx,
				"request failed",
				zap.String("method", info.FullMethod),
				zap.Error(err),
				zap.Duration("duration", duration),
			)
		} else {
			logger.Info(ctx,
				"request completed",
				zap.String("method", info.FullMethod),
				zap.Any("response", resp),
				zap.Duration("duration", duration),
			)
		}

		return resp, err
	}
}
