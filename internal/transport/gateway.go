package transport

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartGateway(ctx context.Context, grpcPort, gatewayPort string) {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "starting grpc gateway...", zap.String("port", gatewayPort))
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := api.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to register order service handler", zap.Error(err))
	}

	if err := http.ListenAndServe(":"+gatewayPort, mux); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to start grpc-gateway", zap.Error(err))
	}
}
