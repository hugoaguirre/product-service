package db

import (
	"context"
	"database/sql"

	"github.com/hugoaguirre/product-service/internal/adapters/db/generated"
	"github.com/hugoaguirre/product-service/internal/domain"
)

type SQLiteRepository struct {
	db      *sql.DB
	queries *generated.Queries
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db:      db,
		queries: generated.New(db),
	}
}

func (r *SQLiteRepository) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	p, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound // or a custom domain error like: domain.ErrNotFound
		}
		return nil, err
	}

	return &domain.Product{
		ID:           p.ID,
		Name:         p.Name,
		PriceInCents: p.PriceInCents,
		Stock:        int32(p.Stock),
	}, nil
}
