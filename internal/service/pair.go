package service

import (
	"slices"
	"strings"

	"exchanger/internal/domain"
	"exchanger/internal/domain/entity"
)

func (s *Service) parsePair(raw string) (entity.Pair, error) {
	base, quote, _ := strings.Cut(raw, "/")
	if !slices.Contains(s.currencies, base) || !slices.Contains(s.currencies, quote) || base == quote {
		return entity.Pair{}, domain.ErrInvalidPair
	}
	return entity.Pair{Base: base, Quote: quote}, nil
}
