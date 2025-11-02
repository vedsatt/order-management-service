package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

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

	orderRepository := repository.NewOrderRepository()

	srv := transport.NewOrderServer(orderRepository)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logger.LoggerInterceptor(ctx)))

	api.RegisterOrderServiceServer(grpcServer, srv)

	go srv.Start(ctx, grpcServer, cfg.GrpcPort)

	go transport.StartGateway(ctx, cfg.GrpcPort, cfg.GatewayPort)

	waitForGracefulShotdown(ctx, grpcServer)
}

func waitForGracefulShotdown(ctx context.Context, grpcServer *grpc.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.GetLoggerFromCtx(ctx).Info(ctx, "shutting down servers...")
	grpcServer.GracefulStop()
	logger.GetLoggerFromCtx(ctx).Info(ctx, "server gracefully stopped")
	fmt.Println("Server Stopped")
}
