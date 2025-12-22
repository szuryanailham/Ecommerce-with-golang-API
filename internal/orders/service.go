package orders

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	repo "github.com/szuryanailham/ecom/internal/adapters/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNotStock = errors.New("product has not enought stock")
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) PlaceOrder(
	ctx context.Context,
	tempOrder createOrderParams,
) (repo.Order, error) {

	// validation payload
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	for _, item := range tempOrder.Items {

		product, err := qtx.FindProductByID(ctx, int64(item.ProductID))
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}

		// atomic stock update
		rows, err := qtx.UpdateQuantityProductByID(ctx, repo.UpdateQuantityProductByIDParams{
			Quantity: item.Quantity,
			ID:       product.ID,
		})
		if err != nil {
			return repo.Order{}, err
		}

		if rows == 0 {
			return repo.Order{}, ErrProductNotStock
		}

		// create order item
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:     order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: product.PriceCents,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}

	// IMPORTANT: check commit error
	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}

	return order, nil
}
