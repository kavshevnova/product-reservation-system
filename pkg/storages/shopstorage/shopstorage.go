package shopstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	"github.com/lib/pq"
	"log"
	"time"
)

type StorageProducts struct {
	db *sqlx.DB
}

func NewShopStorage(dsn string) (*StorageProducts, error) {
	const op = "storages.NewShopStorage"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &StorageProducts{db: db}, nil
}

func (s *StorageProducts) ListProducts(ctx context.Context, limit, offset int32) ([]models.Product, error) {
	const op = "storages.shopstorage.ListProducts"
	const query = "SELECT product_id, name, price, stock FROM products ORDER BY product_id LIMIT $1 OFFSET $2"

	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ProductID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return products, nil
}

func (s *StorageProducts) Product(ctx context.Context, productID int64) (*models.Product, error) {
	const op = "storages.shopstorage.Product"
	const query = "SELECT product_id, name, price, stock FROM products WHERE product_id = $1"

	var product models.Product
	err := s.db.GetContext(ctx, &product, query, productID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &product, nil
}

func (s *StorageProducts) ReserveProduct(ctx context.Context, userID, productID int64, quantity int32) (*models.Order, error) {
	const op = "storages.shopstorage.ReserveProduct"

	log.Printf("Reserving product for user %d, product %d, quantity %d",
		userID, productID, quantity)

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	//Проверяем и блокируем товар
	var price float64
	var stock int32
	err = tx.QueryRowContext(ctx, `SELECT price, stock FROM products WHERE product_id = $1 FOR UPDATE`, productID).Scan(&price, &stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if stock < quantity {
		return nil, models.ErrNotEnoughStock
	}

	//Создаем резервацию
	sum := price * float64(quantity)
	var orderID int64
	now := time.Now()
	err = tx.QueryRowContext(ctx, `INSERT INTO orders (user_id, product_id, quantity, sum, status, time) VALUES ($1, $2, $3, $4, 'reserved', $5) RETURNING order_id`, userID, productID, quantity, sum, now).Scan(&orderID)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, models.ErrOrderAlreadyExists
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//Обновляем остатки
	_, err = tx.ExecContext(ctx, `UPDATE products SET stock = stock - $1 WHERE product_id = $2`, quantity, productID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &models.Order{
		ID:        orderID,
		ProductID: productID,
		UserID:    userID,
		Quantity:  quantity,
		Sum:       float32(sum),
		Status:    "reserved",
		Time:      now,
	}, nil
}

func (s *StorageProducts) ConfirmOrder(ctx context.Context, orderID int64) (*models.Order, error) {
	const op = "storages.shopstorage.ConfirmOrder"
	const query = "UPDATE orders SET status = 'confirmed' WHERE order_id = $1 AND status = 'reserved' RETURNING order_id, user_id, product_id, quantity, sum, status, time"

	var order models.Order
	err := s.db.QueryRowContext(ctx, query, orderID).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Sum, &order.Status, &order.Time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrOrderNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &order, nil
}

func (s *StorageProducts) CancelReservation(ctx context.Context, orderID int64) error {
	const op = "storages.shopstorage.CancelReservation"
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	//Получаем информацию о резервации
	var productID int64
	var quantity int32
	err = tx.QueryRowContext(ctx, `SELECT product_id, quantity FROM orders WHERE order_id = $1 AND status = 'reserved' FOR UPDATE`, orderID).Scan(&productID, &quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrOrderNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	//Возвращаем товар на склад
	_, err = tx.ExecContext(ctx, `UPDATE products SET quantity = quantity + $1 WHERE product_id = $2`, quantity, productID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	//Отменяем резервацию
	_, err = tx.ExecContext(ctx, `UPDATE orders SET status = 'canceled' WHERE order_id = $1`, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *StorageProducts) GetOrderHistory(ctx context.Context, userID int64) ([]models.Order, error) {
	const op = "storages.shopstorage.OrderHistory"
	const query = "SELECT order_id, user_id, product_id, quantity, sum, status, time FROM orders WHERE user_id = $1 ORDER BY time DESC"
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var time sql.NullTime
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Sum, &order.Status, &time); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if time.Valid {
			order.Time = time.Time
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
