package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
)

type OrderService struct {
	mu     sync.RWMutex
	orders map[string]*api.Order
}

func New() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, item string, quantity int32) (string, error) {
	if quantity <= 0 {
		return "", fmt.Errorf("incorrect quantity: %d", quantity)
	}

	id := GenerateOrderID()
	newOrder := &api.Order{
		Id:       id,
		Item:     item,
		Quantity: quantity,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[id] = newOrder

	return id, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*api.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.orders[id]; !ok {
		return &api.Order{}, fmt.Errorf("order with id %s does not exists", id)
	}

	item := s.orders[id]
	return item, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[id]; !ok {
		return &api.Order{}, fmt.Errorf("order with id %s does not exists", id)
	}

	if quantity <= 0 {
		return &api.Order{}, fmt.Errorf("incorrect quantity: %d", quantity)
	}

	order := s.orders[id]
	order.Item = item
	order.Quantity = quantity

	return order, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[id]; !ok {
		return false
	}

	delete(s.orders, id)
	return true
}

func (s *OrderService) ListOrders(ctx context.Context) []*api.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*api.Order, 0, len(s.orders))
	for order := range s.orders {
		orders = append(orders, s.orders[order])
	}

	return orders
}

func GenerateOrderID() string {
	return fmt.Sprintf("order%d", time.Now().UnixNano())
}
