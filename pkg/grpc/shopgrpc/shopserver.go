package shopgrpc

import (
	"context"
	"errors"
	shopv1 "github.com/kavshevnova/product-reservation-system/gen/go/shop"
	"github.com/kavshevnova/product-reservation-system/pkg/domain/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Shop interface {
	ListProducts(ctx context.Context, limit int32, offset int32) ([]models.Product, error)
	GetProductInfo(ctx context.Context, productID int64) (*models.Product, error)
	MakeOrder(ctx context.Context, userID, productID int64, quantity int32) (*models.Order, error)
	GetOrdersHistory(ctx context.Context, userID int64) ([]models.Order, error)
	ConfirmPayment(ctx context.Context, orderID int64, success bool) error
}

type ShopServerAPI struct {
	shopv1.UnimplementedShopServiceServer
	shop Shop
}

func RegisterShopServerAPI(grpcServer *grpc.Server, shop Shop) {
	shopv1.RegisterShopServiceServer(grpcServer, &ShopServerAPI{shop: shop})
}

func (s *ShopServerAPI) ListProducts(ctx context.Context, req *shopv1.ListProductsRequest) (*shopv1.ListProductsResponse, error) {
	if err := ValidateListProducts(req); err != nil {
		return nil, err
	}
	products, err := s.shop.ListProducts(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list products")
	}
	var listProducts []*shopv1.Product
	for _, product := range products {
		listProducts = append(listProducts, &shopv1.Product{
			ProductId: product.ProductID,
			Name:      product.Name,
			Price:     product.Price,
			Stock:     product.Stock,
		})
	}
	return &shopv1.ListProductsResponse{Products: listProducts}, nil
}

func (s *ShopServerAPI) GetProductInfo(ctx context.Context, req *shopv1.GetProductInfoRequest) (*shopv1.GetProductInfoResponse, error) {
	if req.GetProductId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	product, err := s.shop.GetProductInfo(ctx, req.GetProductId())
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to get product")
	}
	return &shopv1.GetProductInfoResponse{
		ProductId: product.ProductID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
	}, nil

}

func (s *ShopServerAPI) MakeOrder(ctx context.Context, req *shopv1.MakeOrderRequest) (*shopv1.MakeOrderResponse, error) {
	if err := ValidateOrderRequest(req); err != nil {
		return nil, err
	}
	order, err := s.shop.MakeOrder(ctx, req.GetUserId(), req.GetProductId(), req.GetQuantity())
	if err != nil {
		switch {
		case errors.Is(err, models.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, "product not found")
		case errors.Is(err, models.ErrNotEnoughStock):
			return &shopv1.MakeOrderResponse{OrderId: order.ID, Status: "Not enough stock"}, nil
		default:
			return nil, status.Error(codes.Internal, "failed to make order")
		}
	}
	return &shopv1.MakeOrderResponse{
		OrderId:    order.ID,
		PaymentURL: order.PaymentURL,
		Status:     order.Status,
	}, nil
}

func (s *ShopServerAPI) GetOrdersHistory(ctx context.Context, req *shopv1.OrdersHistoryRequest) (*shopv1.OrdersHistoryResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	orderHistory, err := s.shop.GetOrdersHistory(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get order history")
	}
	var listOrders []*shopv1.Order
	for _, orders := range orderHistory {
		listOrders = append(listOrders, &shopv1.Order{
			Id:        orders.ID,
			UserId:    orders.UserID,
			ProductId: orders.ProductID,
			Quantity:  orders.Quantity,
			Sum:       orders.Sum,
			OrderTime: orders.Time.Format("2006-01-02 15:04:05.999999999"),
			Status:    orders.Status,
		})
	}
	return &shopv1.OrdersHistoryResponse{Orders: listOrders}, nil
}

func (s *ShopServerAPI) ConfirmPayment(ctx context.Context, req *shopv1.PaymentConfirmation) (*emptypb.Empty, error) {
	if err := s.shop.ConfirmPayment(ctx, req.GetOrderId(), req.GetSuccess()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *ShopServerAPI) mustEmbedUnimplementedShopServiceServer() {}

func ValidateListProducts(request *shopv1.ListProductsRequest) error {
	if request.GetLimit() <= 0 {
		return status.Error(codes.InvalidArgument, "limit must be positive")
	}
	if request.GetOffset() < 0 {
		return status.Error(codes.InvalidArgument, "offset cannot be negative")
	}
	return nil
}

func ValidateOrderRequest(request *shopv1.MakeOrderRequest) error {
	if request.GetProductId() <= 0 {
		return status.Error(codes.InvalidArgument, "product_id is required")
	}
	if request.GetQuantity() <= 0 {
		return status.Error(codes.InvalidArgument, "quantity must be positive")
	}
	if request.GetUserId() <= 0 {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
