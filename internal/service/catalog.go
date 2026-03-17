package service

import (
	"context"

	"github.com/hugoaguirre/product-service/internal/domain"
)

type CatalogService struct {
	repo domain.ProductRepository
}

func NewCatalogService(r domain.ProductRepository) *CatalogService {
	return &CatalogService{repo: r}
}

func (s *CatalogService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return s.repo.GetProduct(ctx, id)
}
