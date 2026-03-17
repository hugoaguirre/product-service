package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/hugoaguirre/product-service/internal/domain"
)

type ProductService interface {
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
}

type ProductJSON struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price_in_cents"`
	Stock int32  `json:"stock"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	svc ProductService
}

func NewHandler(svc ProductService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			h.respondWithError(w, http.StatusBadRequest, "missing product id")
			return
		}

		p, err := h.svc.GetProduct(r.Context(), id)
		if err != nil {
			if errors.Is(err, domain.ErrProductNotFound) {
				h.respondWithError(w, http.StatusNotFound, "product not found")
				return
			}
			// fallback to unexpected errors
			h.respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		resp := ProductJSON{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.PriceInCents,
			Stock: p.Stock,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Default().Printf("unable to encode response: %v\n", err)
		}
	}
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		log.Default().Printf("unable to encode error: %v\n", err)
	}
}
