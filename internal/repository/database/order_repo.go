package database

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	api "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/pkg/api/test"
)

type PostgresCfg struct {
	Host     string `env:"POSTGRES_HOST"     env-default:"postgres"`
	Port     string `env:"POSTGRES_PORT"     env-default:"5432"`
	User     string `env:"POSTGRES_USER"     env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DBName   string `env:"POSTGRES_DB"       env-default:"postgres"`
}

type OrdersDB struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewOrderDB(ctx context.Context, cfg PostgresCfg) (*OrdersDB, error) {
	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	pool, err := pgxpool.New(ctx, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &OrdersDB{
		db:      pool,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (d *OrdersDB) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d *OrdersDB) InsertOrder(ctx context.Context, item string, quantity int32) (string, error) {
	query, args, err := d.builder.Insert("orders").
		Columns("item", "quantity").
		Values(item, quantity).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return "", fmt.Errorf("insert: %w", err)
	}

	var id string
	err = d.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert: %w", err)
	}

	return id, nil
}

func (d *OrdersDB) SelectOrder(ctx context.Context, id string) (*api.Order, error) {
	query, args, err := d.builder.Select(
		"id", "item", "quantity").
		From("orders").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	order := api.Order{}
	err = d.db.QueryRow(ctx, query, args...).Scan(&order.Id, &order.Item, &order.Quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("select: order with id %s does not exists", id)
		}
		return nil, fmt.Errorf("select: %w", err)
	}

	return &order, nil
}

func (d *OrdersDB) UpdateOrder(ctx context.Context, id string, item string, quantity int32) (*api.Order, error) {
	query, args, err := d.builder.Update("orders").
		Set("item", item).
		Set("quantity", quantity).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, item, quantity").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	order := &api.Order{}
	err = d.db.QueryRow(ctx, query, args...).Scan(&order.Id, &order.Item, &order.Quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("select: order with id %s does not exists", id)
		}
		return nil, fmt.Errorf("update: %w", err)
	}

	return order, nil
}

func (d *OrdersDB) DeleteOrder(ctx context.Context, id string) (bool, error) {
	query, args, err := d.builder.Delete("orders").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("delete: %w", err)
	}

	res, err := d.db.Exec(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("delete: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return false, fmt.Errorf("delete: order with id %s does not exists", id)
	}

	return true, nil
}

func (d *OrdersDB) SelectOrdersList(ctx context.Context) ([]*api.Order, error) {
	query, args, err := d.builder.Select(
		"id", "item", "quantity").
		From("orders").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	rows, err := d.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	defer rows.Close()

	orders := make([]*api.Order, 0)
	for rows.Next() {
		order := &api.Order{}
		if err := rows.Scan(&order.Id, &order.Item, &order.Quantity); err != nil {
			return nil, fmt.Errorf("select: %w", err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return orders, nil
}

func (d *OrdersDB) SelectOrdersForCache(ctx context.Context, limit uint64) ([]*api.Order, error) {
	query := "SELECT id, item, quantity FROM orders LIMIT $1"

	rows, err := d.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select orders: %w", err)
	}
	defer rows.Close()

	var orders []*api.Order
	for rows.Next() {
		order := &api.Order{}
		if err := rows.Scan(&order.Id, &order.Item, &order.Quantity); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}
