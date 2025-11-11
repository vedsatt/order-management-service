package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) InsertOrder(ctx context.Context, item string, quantity int32) (string, error) {
	args := m.Called(ctx, item, quantity)
	return args.String(0), args.Error(1)
}

func (m *MockOrderRepository) SelectOrder(ctx context.Context, id string) (*api.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*api.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error) {
	args := m.Called(ctx, id, item, quantity)
	return args.Get(0).(*api.Order), args.Error(1)
}

func (m *MockOrderRepository) DeleteOrder(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderRepository) ListOrders(ctx context.Context) ([]*api.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*api.Order), args.Error(1)
}

func initialize() (*MockOrderRepository, *OrderService, context.Context) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()

	return mockRepo, service, ctx
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockRepo, service, ctx := initialize()

	mockRepo.On("InsertOrder", ctx, "laptop", int32(3)).
		Return("123", nil)

	id, err := service.CreateOrder(ctx, "laptop", 3)

	assert.NoError(t, err)
	assert.Equal(t, "123", id)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_CreateOrder_ValidationError(t *testing.T) {
	mockRepo, service, ctx := initialize()

	_, err := service.CreateOrder(ctx, "", int32(3))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "item cannot be empty")

	_, err = service.CreateOrder(ctx, "laptop", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity must be positive")

	mockRepo.AssertNotCalled(t, "InsertOrder")
}

func TestOrderService_GetOrder_Success(t *testing.T) {
	mockRepo, service, ctx := initialize()

	expected := &api.Order{
		Id:       "12",
		Item:     "bed",
		Quantity: 1,
	}

	mockRepo.On("SelectOrder", ctx, "12").
		Return(expected, nil)

	order, err := service.GetOrder(ctx, "12")

	assert.NoError(t, err)
	assert.Equal(t, expected, order)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetOrder_NotFound(t *testing.T) {
	mockRepo, service, ctx := initialize()

	mockRepo.On("SelectOrder", ctx, "999").
		Return((*api.Order)(nil), fmt.Errorf("not found"))

	order, err := service.GetOrder(ctx, "999")
	assert.Error(t, err)
	assert.Nil(t, order)

	mockRepo.AssertExpectations(t)
}
