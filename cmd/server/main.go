package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/config"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/service"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository/cache"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository/database"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/logger"
)

type App struct {
	GRPCServer    *grpc.Server
	GatewayServer *http.Server
	DB            *database.OrdersDB
	Redis         *cache.OrdersCache
	WG            sync.WaitGroup
}

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

	app := &App{WG: sync.WaitGroup{}}
	app.initializeComponents(ctx, cfg)

	app.gracefulShotdown(ctx)
}

func (a *App) initializeComponents(ctx context.Context, cfg *config.Config) {
	log := logger.GetLoggerFromCtx(ctx)

	var err error
	a.DB, err = database.NewOrderDB(ctx, cfg.PostgresCfg)
	if err != nil {
		log.Fatal(ctx, "database error", zap.Error(err))
	}
	log.Info(ctx, "successfully connected to database")

	a.Redis, err = cache.NewOrdersCache(ctx, cfg.RedisCfg)
	if err != nil {
		log.Fatal(ctx, "redis error", zap.Error(err))
	}
	log.Info(ctx, "succesfully connected to redis")

	orderRepository := repository.NewOrderRepository(a.DB, a.Redis)
	orderService := service.NewOrderService(orderRepository)

	const defaultOrdersLimit = uint64(500)
	go orderRepository.WarmUpCache(ctx, defaultOrdersLimit)

	srv := transport.NewOrderServer(orderService)
	a.GRPCServer = grpc.NewServer(grpc.UnaryInterceptor(logger.LoggerInterceptor(ctx)))
	api.RegisterOrderServiceServer(a.GRPCServer, srv)

	go func() {
		log.Info(ctx, "starting gRPC server...", zap.String("port", cfg.GrpcPort))
		if err = srv.Start(ctx, a.GRPCServer, cfg.GrpcPort); err != nil {
			log.Fatal(ctx, "failed to start gRPC server", zap.Error(err))
		}
	}()

	log.Info(ctx, "starting gRPC gateway...", zap.String("port", cfg.GatewayPort))
	a.GatewayServer, err = transport.StartGateway(ctx, cfg.GrpcPort, cfg.GatewayPort)
	if err != nil {
		log.Fatal(ctx, "failed to start gRPC gateway", zap.Error(err))
	}
}

func (a *App) gracefulShotdown(ctx context.Context) {
	log := logger.GetLoggerFromCtx(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Info(ctx, "shutdown signl received")

	const defaultShutdownTTL = time.Second * 10
	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTTL)
	defer cancel()

	log.Info(ctx, "shutting down gRPC server...")
	a.GRPCServer.GracefulStop()
	log.Info(ctx, "gRPC server stopped successfully")

	log.Info(ctx, "shutting down gRPC gateway...")
	if err := a.GatewayServer.Shutdown(shutdownCtx); err != nil {
		log.Error(ctx, "failed to shutdown gRPC gateway", zap.Error(err))
	} else {
		log.Info(ctx, "gRPC gateway stopped successfully")
	}

	log.Info(ctx, "waiting for background operations...")
	done := make(chan struct{})
	go func() {
		a.WG.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info(ctx, "all background operations completed")
	case <-shutdownCtx.Done():
		log.Info(ctx, "shutdown context done")
	}

	a.closeConnections(shutdownCtx)
	log.Info(ctx, "app shutdown completed")
}

func (a *App) closeConnections(ctx context.Context) {
	if a.Redis != nil {
		zap.L().Info("waiting for cache operations...")
		a.Redis.Wait()
	}

	if a.DB != nil {
		zap.L().Info("closing database connection...")
		a.DB.Close()
		zap.L().Info("database connection closed")
	}

	if a.Redis != nil {
		zap.L().Info("closing Redis connection...")
		a.Redis.Close(ctx)
		zap.L().Info("Redis connection closed")
	}
}
