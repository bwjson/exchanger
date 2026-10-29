package rates

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"exchanger/internal/domain"
	"exchanger/internal/domain/entity"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateUpdate(ctx context.Context, u entity.RateUpdate) error {
	if _, err := r.pool.Exec(ctx, createUpdateStmt, u.ID, u.Pair.Base, u.Pair.Quote, u.Status); err != nil {
		return fmt.Errorf("insert rate update: %w", err)
	}
	return nil
}

func (r *Repository) GetUpdate(ctx context.Context, id uuid.UUID) (entity.RateUpdate, error) {
	u := entity.RateUpdate{ID: id}
	err := r.pool.QueryRow(ctx, getUpdateStmt, id).Scan(
		&u.Pair.Base, &u.Pair.Quote, &u.Status, &u.Price, &u.RateTimestamp, &u.Error, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.RateUpdate{}, domain.ErrNotFound
	}
	if err != nil {
		return entity.RateUpdate{}, fmt.Errorf("select rate update: %w", err)
	}
	return u, nil
}

func (r *Repository) GetLatest(ctx context.Context, pair entity.Pair) (entity.Rate, error) {
	rate := entity.Rate{Pair: pair}
	err := r.pool.QueryRow(ctx, getLatestStmt, pair.Base, pair.Quote).Scan(
		&rate.Price, &rate.RateTimestamp, &rate.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Rate{}, domain.ErrNotFound
	}
	if err != nil {
		return entity.Rate{}, fmt.Errorf("select latest rate: %w", err)
	}
	return rate, nil
}
