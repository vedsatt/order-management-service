package repository

import (
	"context"
	"fmt"

	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository/cache"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository/database"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/logger"
	"go.uber.org/zap"
)

type OrderRepository struct {
	db    *database.OrdersDB
	cache *cache.OrdersCache
}

func NewOrderRepository(db *database.OrdersDB, cache *cache.OrdersCache) *OrderRepository {
	return &OrderRepository{
		db:    db,
		cache: cache,
	}
}

func (r *OrderRepository) WarmUpCache(ctx context.Context, limit uint64) {
	log := logger.GetLoggerFromCtx(ctx)

	orders, err := r.db.SelectOrdersForCache(ctx, limit)
	if err != nil {
		log.Error(ctx, "failed to warm up cache", zap.Error(err))
		return
	}

	for _, order := range orders {
		r.cache.SetOrder(ctx, order)
	}

	log.Info(ctx, "cache warm up completed",
		zap.Int("orders_cached", len(orders)),
		zap.Uint64("limit", limit),
	)
}

func (r *OrderRepository) InsertOrder(ctx context.Context, item string, quantity int32) (string, error) {
	id, err := r.db.InsertOrder(ctx, item, quantity)
	if err != nil {
		return "", fmt.Errorf("database: %w", err)
	}

	order := &api.Order{
		Id:       id,
		Item:     item,
		Quantity: quantity,
	}
	go r.cache.SetOrder(ctx, order)

	return id, nil
}

func (r *OrderRepository) SelectOrder(ctx context.Context, id string) (*api.Order, error) {
	log := logger.GetLoggerFromCtx(ctx)

	order, err := r.cache.GetOrder(ctx, id)
	if err != nil {
		log.Error(ctx, "cache error", zap.Error(err))
	} else if order != nil {
		return order, nil
	}

	order, err = r.db.SelectOrder(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	log.Debug(ctx, "order was found in database", zap.String("id", id))
	return order, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error) {
	order, err := r.db.UpdateOrder(ctx, id, item, quantity)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	r.cache.SetOrder(ctx, order)

	return order, nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id string) (bool, error) {
	success, err := r.db.DeleteOrder(ctx, id)
	if err != nil {
		return success, fmt.Errorf("database: %w", err)
	}

	r.cache.DeleteOrder(ctx, id)

	return success, nil
}

func (r *OrderRepository) ListOrders(ctx context.Context) ([]*api.Order, error) {
	orders, err := r.db.SelectOrdersList(ctx)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	return orders, nil
}
