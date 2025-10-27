package main

import (
	"context"
	"net"

	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/config"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		zap.L().Error("error with config, using default settings", zap.Error(err))
	}

	ctx, err = logger.New(ctx, cfg.Environment)
	if err != nil {
		zap.L().Error("error with creating logger", zap.Error(err))
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "starting server", zap.String("port", cfg.GrpcPort))
	orderRepository := repository.NewOrderRepository()

	srv := transport.NewOrderServer(orderRepository)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logger.LoggerInterceptor(ctx)))

	api.RegisterOrderServiceServer(grpcServer, srv)

	lis, _ := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err = grpcServer.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "error with starting server", zap.Error(err))
	}
}
