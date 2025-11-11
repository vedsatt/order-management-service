package transport

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/logger"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, item string, quantity int32) (string, error) {
	args := m.Called(ctx, item, quantity)
	return args.String(0), args.Error(1)
}

func (m *MockOrderService) GetOrder(ctx context.Context, id string) (*api.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*api.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error) {
	args := m.Called(ctx, id, item, quantity)
	return args.Get(0).(*api.Order), args.Error(1)
}

func (m *MockOrderService) DeleteOrder(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderService) ListOrders(ctx context.Context) ([]*api.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*api.Order), args.Error(1)
}

func TestOrderServer_CreateOrder(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	mockService.On("CreateOrder", mock.Anything, "laptop", int32(2)).
		Return("123", nil)

	ctx, _ := logger.New(context.Background(), "")

	req := &api.CreateOrderRequest{Item: "laptop", Quantity: 2}
	resp, err := server.CreateOrder(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "123", resp.GetId())
	mockService.AssertExpectations(t)
}

func TestOrderServer_CreateOrder_ValidationError(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	mockService.On("CreateOrder", mock.Anything, "", int32(2)).
		Return("", errors.New("item cannot be empty"))
	mockService.On("CreateOrder", mock.Anything, "laptop", int32(0)).
		Return("", errors.New("quantity must be positive"))

	ctx, _ := logger.New(context.Background(), "")

	req := &api.CreateOrderRequest{Item: "", Quantity: 2}
	resp, err := server.CreateOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "item cannot be empty")

	req = &api.CreateOrderRequest{Item: "laptop", Quantity: 0}
	resp, err = server.CreateOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "quantity must be positive")

	mockService.AssertExpectations(t)
}

func TestOrderServer_GetOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	expectedOrder := &api.Order{Id: "123", Item: "laptop", Quantity: 2}
	mockService.On("GetOrder", mock.Anything, "123").Return(expectedOrder, nil)

	ctx, _ := logger.New(context.Background(), "")

	req := &api.GetOrderRequest{Id: "123"}
	resp, err := server.GetOrder(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, resp.GetOrder())
	mockService.AssertExpectations(t)
}

func TestOrderServer_GetOrder_NotFound(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	mockService.On("GetOrder", mock.Anything, "999").
		Return((*api.Order)(nil), errors.New("not found"))

	ctx, _ := logger.New(context.Background(), "")

	req := &api.GetOrderRequest{Id: "999"}
	resp, err := server.GetOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestOrderServer_UpdateOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	expectedOrder := &api.Order{Id: "123", Item: "updated-laptop", Quantity: 3}
	mockService.On("UpdateOrder", mock.Anything, "123", "updated-laptop", int32(3)).
		Return(expectedOrder, nil)

	ctx, _ := logger.New(context.Background(), "")

	req := &api.UpdateOrderRequest{Id: "123", Item: "updated-laptop", Quantity: 3}
	resp, err := server.UpdateOrder(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, resp.GetOrder())
	mockService.AssertExpectations(t)
}

func TestOrderServer_UpdateOrder_ValidationError(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	mockService.On("UpdateOrder", mock.Anything, "123", "", int32(3)).
		Return((*api.Order)(nil), errors.New("item cannot be empty"))
	mockService.On("UpdateOrder", mock.Anything, "123", "laptop", int32(0)).
		Return((*api.Order)(nil), errors.New("quantity must be positive"))

	ctx, _ := logger.New(context.Background(), "")

	req := &api.UpdateOrderRequest{Id: "123", Item: "", Quantity: 3}
	resp, err := server.UpdateOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "item cannot be empty")

	req = &api.UpdateOrderRequest{Id: "123", Item: "laptop", Quantity: 0}
	resp, err = server.UpdateOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "quantity must be positive")

	mockService.AssertExpectations(t)
}

func TestOrderServer_DeleteOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	mockService.On("DeleteOrder", mock.Anything, "123").Return(true, nil)

	ctx, _ := logger.New(context.Background(), "")

	req := &api.DeleteOrderRequest{Id: "123"}
	resp, err := server.DeleteOrder(ctx, req)

	assert.NoError(t, err)
	assert.True(t, resp.GetSuccess())
	mockService.AssertExpectations(t)
}

func TestOrderServer_DeleteOrder_NotFound(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	mockService.On("DeleteOrder", mock.Anything, "999").Return(false, errors.New("not found"))

	ctx, _ := logger.New(context.Background(), "")

	req := &api.DeleteOrderRequest{Id: "999"}
	resp, err := server.DeleteOrder(ctx, req)

	assert.Error(t, err)
	assert.False(t, resp.GetSuccess())
	mockService.AssertExpectations(t)
}

func TestOrderServer_ListOrders_Success(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	expectedOrders := []*api.Order{
		{Id: "1", Item: "laptop", Quantity: 2},
		{Id: "2", Item: "mouse", Quantity: 5},
	}
	mockService.On("ListOrders", mock.Anything).Return(expectedOrders, nil)

	ctx, _ := logger.New(context.Background(), "")

	req := &api.ListOrdersRequest{}
	resp, err := server.ListOrders(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, resp.GetOrders())
	mockService.AssertExpectations(t)
}

func TestOrderServer_ListOrders_Empty(t *testing.T) {
	mockService := new(MockOrderService)
	server := NewOrderServer(mockService)

	emptyOrders := []*api.Order{}
	mockService.On("ListOrders", mock.Anything).Return(emptyOrders, nil)

	ctx, _ := logger.New(context.Background(), "")

	req := &api.ListOrdersRequest{}
	resp, err := server.ListOrders(ctx, req)

	assert.NoError(t, err)
	assert.Empty(t, resp.GetOrders())
	mockService.AssertExpectations(t)
}
