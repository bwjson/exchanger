package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"exchanger/internal/domain/entity"
	"exchanger/internal/domain/enum"
)

type Repository interface {
	CreateUpdate(ctx context.Context, u entity.RateUpdate) error
	GetUpdate(ctx context.Context, id uuid.UUID) (entity.RateUpdate, error)
	GetLatest(ctx context.Context, pair entity.Pair) (entity.Rate, error)
}

type Service struct {
	repo       Repository
	currencies []string
}

func New(repo Repository, currencies []string) *Service {
	return &Service{repo: repo, currencies: currencies}
}

func (s *Service) CreateUpdate(ctx context.Context, rawPair string) (uuid.UUID, error) {
	pair, err := s.parsePair(rawPair)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate update id: %w", err)
	}

	update := entity.RateUpdate{ID: id, Pair: pair, Status: enum.UpdateStatusPending}
	if err := s.repo.CreateUpdate(ctx, update); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Service) GetUpdate(ctx context.Context, id uuid.UUID) (entity.RateUpdate, error) {
	return s.repo.GetUpdate(ctx, id)
}

func (s *Service) GetLatest(ctx context.Context, rawPair string) (entity.Rate, error) {
	pair, err := s.parsePair(rawPair)
	if err != nil {
		return entity.Rate{}, err
	}
	return s.repo.GetLatest(ctx, pair)
}
