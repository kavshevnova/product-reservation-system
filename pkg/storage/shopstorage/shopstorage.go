package shopstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	"github.com/lib/pq"
	"time"
)

type StorageProducts struct {
	db *sql.DB
}

func NewShopStorage(db *sql.DB) (*StorageProducts, error) {
	return &StorageProducts{db: db}, nil
}

func (s *StorageProducts) ListProducts(ctx context.Context, limit, offset int32) ([]models.Product, error) {
	const op = "storage.shopstorage.ListProducts"
	const query = "SELECT ProductID, Name, Price, Stock FROM products ORDER BY productID LIMIT $1 OFFSET $2"

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
	const op = "storage.shopstorage.Product"
	const query = "SELECT ProductID, Name, Price, Stock FROM products WHERE ProductID = $1"

	var product models.Product
	err := s.db.QueryRowContext(ctx, query, productID).Scan(&product.ProductID, &product.Name, &product.Price, &product.Stock)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &product, nil
}

func (s *StorageProducts) ReserveProduct(ctx context.Context, userID, productID int64, quantity int32) (*models.Order, error) {
	const op = "storage.shopstorage.ReserveProduct"
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	//Проверяем и блокируем товар
	var price float64
	var stock int32
	err = tx.QueryRowContext(ctx, `SELECT price, stock FROM products WHERE id = $1 FOR UPDATE`, productID).Scan(&price, &stock)
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
	err = tx.QueryRowContext(ctx, `INSERT INTO orders (userID, productID, quantity, Sum, Status, Time) VALUES ($1, $2, $3, $4, 'reserved', $5) RETURNING ID`, userID, productID, quantity, sum, time.Now()).Scan(&orderID)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, models.ErrOrderAlreadyExists
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//Обновляем остатки
	_, err = tx.ExecContext(ctx, `UPDATE products SET stock = stock - $1 WHERE ID = $2`, quantity, productID)
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
		Time:      time.Now(),
	}, nil
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
