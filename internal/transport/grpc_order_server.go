package transport

import (
	"context"
	"strings"

	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	api.UnimplementedOrderServiceServer
	service *repository.OrderRepository
}

func NewOrderServer(srv *repository.OrderRepository) *OrderServer {
	return &OrderServer{
		service: srv,
	}
}

func (o *OrderServer) CreateOrder(ctx context.Context, in *api.CreateOrderRequest) (*api.CreateOrderResponse, error) {
	id, err := o.service.CreateOrder(ctx, in.Item, in.Quantity)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp := &api.CreateOrderResponse{
		Id: id,
	}
	return resp, nil
}

func (o *OrderServer) GetOrder(ctx context.Context, in *api.GetOrderRequest) (*api.GetOrderResponse, error) {
	order, err := o.service.GetOrder(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	resp := &api.GetOrderResponse{
		Order: order,
	}
	return resp, nil
}

func (o *OrderServer) UpdateOrder(ctx context.Context, in *api.UpdateOrderRequest) (*api.UpdateOrderResponse, error) {
	updOrder, err := o.service.UpdateOrder(ctx, in.Id, in.Item, in.Quantity)
	if err != nil {
		if strings.Contains(err.Error(), "does not exists") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp := &api.UpdateOrderResponse{
		Order: updOrder,
	}
	return resp, nil
}

func (o *OrderServer) DeleteOrder(ctx context.Context, in *api.DeleteOrderRequest) (*api.DeleteOrderResponse, error) {
	success, err := o.service.DeleteOrder(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &api.DeleteOrderResponse{
		Success: success,
	}
	return resp, nil
}

func (o *OrderServer) ListOrders(ctx context.Context, in *api.ListOrdersRequest) (*api.ListOrdersResponse, error) {
	orders, err := o.service.ListOrders(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &api.ListOrdersResponse{
		Orders: orders,
	}
	return resp, nil
}
