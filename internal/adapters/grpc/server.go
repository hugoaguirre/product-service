package grpc

import (
	"context"
	"errors"

	"github.com/hugoaguirre/product-service/internal/domain"
	"github.com/hugoaguirre/product-service/internal/service"
	"github.com/hugoaguirre/product-service/pkg/productapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcAdapter struct {
	productapi.UnimplementedProductServiceServer
	service *service.CatalogService
}

func New(svc *service.CatalogService) *GrpcAdapter {
	return &GrpcAdapter{service: svc}
}

// GetProduct returns a product information
func (a *GrpcAdapter) GetProduct(ctx context.Context, r *productapi.ProductRequest) (*productapi.ProductResponse, error) {
	p, err := a.service.GetProduct(ctx, r.ProductId)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &productapi.ProductResponse{
		Name:         p.Name,
		PriceInCents: p.PriceInCents,
		Stock:        p.Stock,
	}, nil
}
