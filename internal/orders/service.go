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

func (s *svc) PlaceOrder (ctx context.Context, tempOrder createOrderParams )(repo.Order, error){
	// validation payload
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("costumer ID is required")
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
	// created an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}
	// create for the product if exist
	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductByID(ctx,int64(item.ProductID))
		if err != nil {
			return repo.Order{},ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{},ErrProductNotStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
		OrderID:order.ID,
		ProductID: item.ProductID,
		Quantity: item.Quantity,
		PriceCents: product.PriceCents,
	})
	if err != nil  {
		return repo.Order{}, err
	}

	// Challenge : update the product stock quantity
	}
	
	tx.Commit(ctx)
	// create  order item
	return order, nil
}