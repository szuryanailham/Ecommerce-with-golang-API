package products

import (
	"context"

	repo "github.com/szuryanailham/ecom/internal/adapters/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, id int32)(repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier)Service{
	return &svc{repo:repo}
}

func (s *svc)ListProducts(ctx context.Context)([]repo.Product, error){
	return s.repo.ListProducts(ctx)
}

func (s *svc)FindProductByID(ctx context.Context, id int32) (repo.Product, error) {
    return s.repo.FindProductByID(ctx, int64(id))
}