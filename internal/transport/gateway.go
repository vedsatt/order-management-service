package transport

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartGateway(ctx context.Context, grpcPort, gatewayPort string) (*http.Server, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := api.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		return nil, fmt.Errorf("failed to register order service handler: %w", err)
	}

	const defaultGatewayTimeout = 5 * time.Second
	server := &http.Server{
		Addr:              ":" + gatewayPort,
		Handler:           mux,
		ReadHeaderTimeout: defaultGatewayTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start gateway: %w", err))
		}
	}()

	return server, nil
}
