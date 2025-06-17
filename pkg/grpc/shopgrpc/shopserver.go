package shopgrpc

import (
	"context"
	shopv1 "github.com/kavshevnova/product-reservation-system/gen/go/shop"
	"reflect"
)

type Shop interface {
}

type ShopServerAPI struct {
	shopv1.UnimplementedShopServiceServer
	shop Shop
}

func (s *ShopServerAPI) ListProducts(ctx context.Context, req *shopv1.ListProductsRequest) (*shopv1.ListProductsResponse, error) {

}

func (s *ShopServerAPI) GetProductInfo(ctx context.Context, req *shopv1.GetProductInfoRequest) (*shopv1.GetProductInfoResponse, error) {

}

func (s *ShopServerAPI) MakeOrder(ctx context.Context, req *shopv1.OrderRequest) (*shopv1.OrderResponse, error) {

}

func (s *ShopServerAPI) GetOrdersHistory(ctx context.Context, req *shopv1.OrdersHistoryRequest) (*shopv1.OrdersHistoryResponse, error) {

}
