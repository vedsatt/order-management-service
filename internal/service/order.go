package service

import (
	"context"
	"fmt"

	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, item string, quantity int32) (string, error)
	SelectOrder(ctx context.Context, id string) (*api.Order, error)
	UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error)
	DeleteOrder(ctx context.Context, id string) (bool, error)
	ListOrders(ctx context.Context) ([]*api.Order, error)
}

type OrderService struct {
	repository OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		repository: repo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, item string, quantity int32) (string, error) {
	if item == "" {
		return "", fmt.Errorf("item cannot be empty")
	}

	if quantity <= 0 {
		return "", fmt.Errorf("quantity must be positive")
	}

	id, err := s.repository.InsertOrder(ctx, item, quantity)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*api.Order, error) {
	order, err := s.repository.SelectOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error) {
	if item == "" {
		return nil, fmt.Errorf("item cannot be empty")
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	order, err := s.repository.UpdateOrder(ctx, id, item, quantity)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) (bool, error) {
	success, err := s.repository.DeleteOrder(ctx, id)
	if err != nil {
		return success, err
	}

	return success, nil
}

func (s *OrderService) ListOrders(ctx context.Context) ([]*api.Order, error) {
	orders, err := s.repository.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
