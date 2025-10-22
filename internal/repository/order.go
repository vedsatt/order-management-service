package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*api.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]*api.Order),
	}
}

func (s *OrderRepository) CreateOrder(ctx context.Context, item string, quantity int32) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if item == "" {
		return "", fmt.Errorf("item cannot be empty")
	}

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

func (s *OrderRepository) GetOrder(ctx context.Context, id string) (*api.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.orders[id]; !ok {
		return nil, fmt.Errorf("order with id %s does not exists", id)
	}

	item := s.orders[id]
	return item, nil
}

func (s *OrderRepository) UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[id]; !ok {
		return nil, fmt.Errorf("order with id %s does not exists", id)
	}

	if item == "" {
		return nil, fmt.Errorf("item cannot be empty")
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("incorrect quantity: %d", quantity)
	}

	order := &api.Order{
		Id:       id,
		Item:     item,
		Quantity: quantity,
	}
	s.orders[id] = order

	return order, nil
}

func (s *OrderRepository) DeleteOrder(ctx context.Context, id string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[id]; !ok {
		return false, nil
	}

	delete(s.orders, id)
	return true, nil
}

func (s *OrderRepository) ListOrders(ctx context.Context) ([]*api.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*api.Order, 0, len(s.orders))
	for order := range s.orders {
		orders = append(orders, s.orders[order])
	}

	return orders, nil
}

func GenerateOrderID() string {
	return fmt.Sprintf("order%d", time.Now().UnixNano())
}
