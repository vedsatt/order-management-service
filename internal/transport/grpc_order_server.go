package transport

import (
	"context"
	"fmt"
	"net"
	"strings"

	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService interface {
	CreateOrder(ctx context.Context, item string, quantity int32) (string, error)
	GetOrder(ctx context.Context, id string) (*api.Order, error)
	UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error)
	DeleteOrder(ctx context.Context, id string) (bool, error)
	ListOrders(ctx context.Context) ([]*api.Order, error)
}

type OrderServer struct {
	api.UnimplementedOrderServiceServer

	service OrderService
}

func NewOrderServer(srv OrderService) *OrderServer {
	return &OrderServer{
		service: srv,
	}
}

func (s *OrderServer) Start(ctx context.Context, grpcServer *grpc.Server, port string) error {
	lc := &net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *OrderServer) CreateOrder(ctx context.Context, in *api.CreateOrderRequest) (*api.CreateOrderResponse, error) {
	log := logger.GetLoggerFromCtx(ctx)

	log.Debug(ctx, "CreateOrder raw request",
		zap.Any("raw request", in),
	)

	log.Info(ctx, "CreateOrder started",
		zap.String("item", in.GetItem()),
		zap.Int32("quantity", in.GetQuantity()),
	)

	id, err := s.service.CreateOrder(ctx, in.GetItem(), in.GetQuantity())
	if err != nil {
		log.Error(ctx, "CreateOrder failed",
			zap.String("item", in.GetItem()),
			zap.Int32("quantity", in.GetQuantity()),
			zap.Error(err),
		)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Info(ctx, "CreateOrder completed",
		zap.String("order_id", id),
		zap.String("item", in.GetItem()),
		zap.Int32("quantity", in.GetQuantity()),
	)

	resp := &api.CreateOrderResponse{
		Id: id,
	}

	log.Debug(ctx, "CreateOrder raw response",
		zap.Any("raw response", resp),
	)

	return resp, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, in *api.GetOrderRequest) (*api.GetOrderResponse, error) {
	log := logger.GetLoggerFromCtx(ctx)

	log.Debug(ctx, "GetOrder raw request",
		zap.Any("raw request", in),
	)

	log.Info(ctx, "GetOrder started",
		zap.String("order_id", in.GetId()),
	)

	order, err := s.service.GetOrder(ctx, in.GetId())
	if err != nil {
		log.Warn(ctx, "GetOrder not found",
			zap.String("order_id", in.GetId()),
			zap.Error(err),
		)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	log.Info(ctx, "GetOrder completed",
		zap.String("order_id", in.GetId()),
		zap.String("item", order.GetItem()),
		zap.Int32("quantity", order.GetQuantity()),
	)

	resp := &api.GetOrderResponse{
		Order: order,
	}

	log.Debug(ctx, "GetOrder raw response",
		zap.Any("raw response", resp),
	)

	return resp, nil
}

func (s *OrderServer) UpdateOrder(ctx context.Context, in *api.UpdateOrderRequest) (*api.UpdateOrderResponse, error) {
	log := logger.GetLoggerFromCtx(ctx)

	log.Debug(ctx, "UpdateOrder raw request",
		zap.Any("raw request", in),
	)

	log.Info(ctx, "UpdateOrder started",
		zap.String("order_id", in.GetId()),
		zap.String("new_item", in.GetItem()),
		zap.Int32("new_quantity", in.GetQuantity()),
	)

	updOrder, err := s.service.UpdateOrder(ctx, in.GetId(), in.GetItem(), in.GetQuantity())
	if err != nil {
		if strings.Contains(err.Error(), "does not exists") {
			log.Warn(ctx, "UpdateOrder not found",
				zap.String("order_id", in.GetId()),
				zap.Error(err),
			)
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.Error(ctx, "UpdateOrder failed",
			zap.String("order_id", in.GetId()),
			zap.Error(err),
		)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Info(ctx, "UpdateOrder completed",
		zap.String("order_id", in.GetId()),
		zap.String("item", updOrder.GetItem()),
		zap.Int32("quantity", updOrder.GetQuantity()),
	)

	resp := &api.UpdateOrderResponse{
		Order: updOrder,
	}

	log.Debug(ctx, "UpdateOrder raw response",
		zap.Any("raw response", resp),
	)

	return resp, nil
}

func (s *OrderServer) DeleteOrder(ctx context.Context, in *api.DeleteOrderRequest) (*api.DeleteOrderResponse, error) {
	log := logger.GetLoggerFromCtx(ctx)

	log.Debug(ctx, "DeleteOrder raw request",
		zap.Any("raw request", in),
	)

	log.Info(ctx, "DeleteOrder started",
		zap.String("order_id", in.GetId()),
	)

	success, err := s.service.DeleteOrder(ctx, in.GetId())
	if err != nil {
		log.Error(ctx, "DeleteOrder failed",
			zap.String("order_id", in.GetId()),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if success {
		log.Info(ctx, "DeleteOrder completed",
			zap.String("order_id", in.GetId()),
			zap.Bool("success", success),
		)
	} else {
		log.Warn(ctx, "DeleteOrder not found",
			zap.String("order_id", in.GetId()),
			zap.Bool("success", success),
		)
	}

	resp := &api.DeleteOrderResponse{
		Success: success,
	}

	log.Debug(ctx, "DeleteOrder raw response",
		zap.Any("raw response", resp),
	)

	return resp, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, in *api.ListOrdersRequest) (*api.ListOrdersResponse, error) {
	log := logger.GetLoggerFromCtx(ctx)

	log.Debug(ctx, "ListOrder raw request",
		zap.Any("raw request", in),
	)

	log.Info(ctx, "ListOrders started")

	orders, err := s.service.ListOrders(ctx)
	if err != nil {
		log.Error(ctx, "ListOrders failed",
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info(ctx, "ListOrders completed",
		zap.Int("orders_count", len(orders)),
	)

	resp := &api.ListOrdersResponse{
		Orders: orders,
	}

	log.Debug(ctx, "ListOrder raw response",
		zap.Any("raw response", resp),
	)

	return resp, nil
}
