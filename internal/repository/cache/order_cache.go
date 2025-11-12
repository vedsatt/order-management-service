package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/logger"
	"go.uber.org/zap"
)

type RedisCfg struct {
	Host     string `env:"REDIS_HOST"     env-default:"redis"`
	Port     string `env:"REDIS_PORT"     env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:"redis"`
}

type OrdersCache struct {
	redisClient *redis.Client
	wg          sync.WaitGroup
}

var (
	ErrOrderNotFound = errors.New("order not found in cache")
)

func NewOrdersCache(ctx context.Context, cfg RedisCfg) (*OrdersCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &OrdersCache{
		redisClient: client,
		wg:          sync.WaitGroup{},
	}, nil
}

func (c *OrdersCache) Wait() {
	c.wg.Wait()
}

func (c *OrdersCache) Close(ctx context.Context) {
	if c.redisClient != nil {
		if err := c.redisClient.Close(); err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "error closing Redis connection", zap.Error(err))
		}
	}
}

func (c *OrdersCache) SetOrder(ctx context.Context, order *api.Order) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		bgCtx := context.Background()
		log := logger.GetLoggerFromCtx(ctx)

		log.Debug(ctx, "SetOrder - marshaling order",
			zap.String("id", order.GetId()),
			zap.String("item", order.GetItem()),
			zap.Int32("quantity", order.GetQuantity()),
		)

		data, err := json.Marshal(order)
		if err != nil {
			log.Error(ctx, "failed to marshal order", zap.Error(err), zap.String("id", order.GetId()))
			return
		}

		log.Debug(ctx, "SetOrder - JSON data",
			zap.String("json", string(data)),
			zap.String("redis_key", order.GetId()),
		)

		const defaultTTL = time.Minute * 30
		err = c.redisClient.Set(bgCtx, order.GetId(), data, defaultTTL).Err()
		if err != nil {
			log.Error(ctx, "failed to set order to redis", zap.Error(err), zap.String("id", order.GetId()))
			return
		}

		log.Debug(ctx, "order successfully set to redis", zap.String("id", order.GetId()))
	}()
}

func (c *OrdersCache) GetOrder(ctx context.Context, id string) (*api.Order, error) {
	log := logger.GetLoggerFromCtx(ctx)

	log.Debug(ctx, "GetOrder - searching in Redis", zap.String("redis_key", id))

	val, err := c.redisClient.Get(ctx, id).Bytes()
	if errors.Is(err, redis.Nil) {
		log.Debug(ctx, "GetOrder - not found in Redis", zap.String("redis_key", id))
		return nil, ErrOrderNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error with cache: %w", err)
	}

	log.Debug(ctx, "GetOrder - raw data from Redis",
		zap.String("redis_key", id),
		zap.String("raw_json", string(val)),
	)

	var order api.Order
	if err = json.Unmarshal(val, &order); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	log.Debug(ctx, "GetOrder - unmarshaled order",
		zap.String("id", order.GetId()),
		zap.String("item", order.GetItem()),
		zap.Int32("quantity", order.GetQuantity()),
	)

	return &order, nil
}

func (c *OrdersCache) DeleteOrder(ctx context.Context, id string) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()
		bgCtx := context.Background()

		log := logger.GetLoggerFromCtx(ctx)

		err := c.redisClient.Del(bgCtx, id).Err()
		if err != nil {
			log.Error(ctx, "failed to delete order from redis", zap.Error(err), zap.String("id", id))
		}

		log.Debug(ctx, "successfully deleted order from redis", zap.String("id", id))
	}()
}
