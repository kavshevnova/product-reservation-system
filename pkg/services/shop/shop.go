package shop

import (
	"context"
	"errors"
	"fmt"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
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
	GetOrderHistory(ctx context.Context, userID int64) ([]models.Order, error)
}

type InventoryManager interface {
	ReserveProduct(ctx context.Context, userID, productID int64, quantity int32) (*models.Order, error)
	CancelReservation(ctx context.Context, orderID int64) error
	ConfirmOrder(ctx context.Context, orderID int64) (*models.Order, error)
}

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
		if errors.Is(err, models.ErrProductNotFound) {
			s.log.Warn("Product not found", slog.StringValue(err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("GetProduct failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Get Product done")
	return product, nil
}

func (s *Shop) GetOrdersHistory(ctx context.Context, userID int64) ([]models.Order, error) {
	const op = "shop.OrderHistory"

	log := s.log.With(slog.String("operation", op), slog.String("userID", strconv.Itoa(int(userID))))
	log.Info("Starting Get OrderHistory")

	orders, err := s.storage.GetOrderHistory(ctx, userID)
	if err != nil {
		log.Error("GetOrderHistory failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info(("Get OrderHistory done"), slog.Int("count", len(orders)))
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

	//Проверяем наличие товара

	product, err := s.storage.Product(ctx, productID)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			log.Error("Product not found", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, models.ErrProductNotFound)
		}
		log.Error("Failed to get product info", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if product.Stock < quantity {
		log.Error("Not enough stock", slog.String("Available", strconv.Itoa(int(product.Stock))))
		return nil, fmt.Errorf("%s: %w", op, models.ErrNotEnoughStock)
	}

	log.Info("Product in stock")

	//Резервируем товар

	order, err := s.inventory.ReserveProduct(ctx, userID, productID, quantity)
	if err != nil {
		log.Error("Failed to reserve product", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Reserve Product done", slog.String("productID", strconv.Itoa(int(order.ID))))

	//Возвращаем заказ в статусе "ожидает оплаты"
	return &models.Order{
		ID:         order.ID,
		Status:     "waiting_payment",
		PaymentURL: s.generatePaymentURL(order.ID),
	}, nil
}

func (s *Shop) generatePaymentURL(orderID int64) string {
	return fmt.Sprintf("https://pay.example.com?order_id=%d", orderID)
}

func (s *Shop) ConfirmPayment(ctx context.Context, orderID int64, success bool) error {
	const op = "services.shop.ConfirmPayment"
	log := s.log.With(
		slog.String("op", op),
		slog.Int64("order_id", orderID),
		slog.Bool("success", success),
	)

	if success {
		// Подтверждаем заказ
		_, err := s.inventory.ConfirmOrder(ctx, orderID)
		if err != nil {
			log.Error("failed to confirm order", slog.String("error", err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}
		log.Info("payment confirmed")
	} else {
		// Отменяем резервацию
		err := s.inventory.CancelReservation(ctx, orderID)
		if err != nil {
			log.Error("failed to cancel reservation", slog.String("error", err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}
		log.Info("payment failed, reservation canceled")
	}

	return nil
}
