package rates

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"exchanger/internal/domain"
	"exchanger/internal/domain/entity"
	"exchanger/internal/domain/enum"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateUpdate(ctx context.Context, u entity.RateUpdate) (entity.RateUpdate, error) {
	err := r.pool.QueryRow(ctx, createUpdateStmt, u.ID, u.Pair.Base, u.Pair.Quote, u.Status, u.IdempotencyKey).Scan(&u.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = r.pool.QueryRow(ctx, getUpdateByKeyStmt, u.IdempotencyKey).Scan(&u.ID, &u.Pair.Base, &u.Pair.Quote)
	}
	if err != nil {
		return entity.RateUpdate{}, fmt.Errorf("insert rate update: %w", err)
	}
	return u, nil
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

func (r *Repository) ProcessPending(ctx context.Context, limit int, resolve func([]entity.RateUpdate)) error {
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		updates, err := selectPending(ctx, tx, limit)
		if err != nil || len(updates) == 0 {
			return err
		}
		resolve(updates)
		return saveUpdates(ctx, tx, updates)
	})
}

func selectPending(ctx context.Context, tx pgx.Tx, limit int) ([]entity.RateUpdate, error) {
	rows, err := tx.Query(ctx, selectPendingStmt, limit)
	if err != nil {
		return nil, fmt.Errorf("select pending updates: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.RateUpdate, error) {
		var u entity.RateUpdate
		err := row.Scan(&u.ID, &u.Pair.Base, &u.Pair.Quote, &u.Attempts, &u.NextAttemptAt)
		return u, err
	})
}

func saveUpdates(ctx context.Context, tx pgx.Tx, updates []entity.RateUpdate) error {
	slices.SortFunc(updates, func(a, b entity.RateUpdate) int {
		return strings.Compare(a.Pair.String(), b.Pair.String())
	})

	batch := &pgx.Batch{}
	for _, u := range updates {
		batch.Queue(saveUpdateStmt, u.ID, u.Status, u.Price, u.RateTimestamp, u.Error, u.Attempts, u.NextAttemptAt)
		if u.Status == enum.UpdateStatusDone {
			batch.Queue(upsertRateStmt, u.Pair.Base, u.Pair.Quote, u.Price, u.RateTimestamp)
		}
	}
	if err := tx.SendBatch(ctx, batch).Close(); err != nil {
		return fmt.Errorf("save updates: %w", err)
	}
	return nil
}
