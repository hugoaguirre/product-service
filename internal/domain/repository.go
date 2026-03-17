package domain

import "context"

// Product is part of the domain model, make sure no tags are added in its fields
type Product struct {
	ID           string
	Name         string
	PriceInCents int64
	Stock        int32
}

type ProductRepository interface {
	// GetProduct retrieves a product from the repository
	GetProduct(ctx context.Context, id string) (*Product, error)

	// Save stores a product into the repository
	// Save(ctx context.Context, p Product) error
}
