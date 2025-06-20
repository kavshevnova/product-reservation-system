package shop

import (
	"context"
	"errors"
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	storageerrors "github.com/kavshevnova/product-reservation-system/pkg/storage"
	"log/slog"
	"strconv"
)

type Shop struct {
	log       *slog.Logger
	storage   ProductStorage
	inventory InventoryManager
}

type ProductStorage interface {
	ListProducts(ctx context.Context, limit, offset int32) ([]models.Product, error)
	Product(ctx context.Context, productID int64) (*models.Product, error)
	OrderHistory(ctx context.Context, userID int64) ([]models.Order, error)
}

type InventoryManager interface {
	BuyProduct(ctx context.Context, userID, productID int64, quantity int32) (*models.Order, error)
	ReserveProduct(ctx context.Context, userID, productID int64, quantity int32) (*models.Order, error)
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrNotEnoughStock  = errors.New("not enough stock")
)

func New(log *slog.Logger, storage ProductStorage, inventory InventoryManager) *Shop {
	return &Shop{
		log:       log,
		storage:   storage,
		inventory: inventory,
	}
}

func (s *Shop) ListProducts(ctx context.Context, limit, offset int32) ([]models.Product, error) {
	const op = "shop.ListProducts"

	log := s.log.With(
		slog.String("operation", op),
		slog.String("limit", strconv.Itoa(int(limit))),
		slog.String("offset", strconv.Itoa(int(offset))),
	)

	log.Info("Starting list Products")

	products, err := s.storage.ListProducts(ctx, limit, offset)
	if err != nil {
		log.Error("ListProducts failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("ListProducts done")
	return products, nil
}

func (s *Shop) GetProductInfo(ctx context.Context, productID int64) (*models.Product, error) {
	const op = "shop.Product"

	log := s.log.With(
		slog.String("operation", op),
		slog.String("productID", strconv.Itoa(int(productID))),
	)

	log.Info("Starting  Get Product")

	product, err := s.storage.Product(ctx, productID)
	if err != nil {
		if errors.Is(err, storageerrors.ErrProductNotFound) {
			s.log.Warn("Product not found", slog.StringValue(err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("GetProduct failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Get Product done")
	return product, nil
}

func (s *Shop) GetOrderHistory(ctx context.Context, userID int64) ([]models.Order, error) {
	const op = "shop.OrderHistory"

	log := s.log.With(slog.String("operation", op), slog.String("userID", strconv.Itoa(int(userID))))
	log.Info("Starting Get OrderHistory")

	orders, err := s.storage.OrderHistory(ctx, userID)
	if err != nil {
		log.Error("GetOrderHistory failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info(("Get OrderHistory done"), slog.Int("count", len(orders))
	return orders, nil
}

func (s *Shop) MakeOrder(ctx context.Context, userID, productID int64, quantity int32) (*models.Order, error) {
	const op = "shop.MakeOrder"

	log := s.log.With(
		slog.String("operation", op),
		slog.String("productID", strconv.Itoa(int(userID))),
		slog.String("productID", strconv.Itoa(int(productID))),
		slog.String("quantity", strconv.Itoa(int(quantity))),
		)
	log.Info("Starting Buy Product")

	product, err := s.GetProductInfo(ctx, productID)
	if product.Stock == 0 {
		log.Error("Product stock is zero", slog.String("productID", strconv.Itoa(int(productID)))
		return nil, fmt.Errorf("%s: %w", op, ErrNotEnoughStock)
	}
	if err != nil {
		log.Error("GetProduct failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Product in stock")



}
