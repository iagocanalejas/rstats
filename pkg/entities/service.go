package entities

import (
	"context"

	"github.com/iagocanalejas/regatas/internal/db"
)

type Service struct {
	db *db.Queries
}

func NewService(d *db.Queries) *Service {
	return &Service{
		db: d,
	}
}

func (s *Service) GetClubs(ctx context.Context) ([]db.Entity, error) {
	return s.db.GetClubs(ctx)
}

func (s *Service) GetEntities(ctx context.Context) ([]db.Entity, error) {
	return s.db.GetEntities(ctx)
}

func (s *Service) GetLeagues(ctx context.Context) ([]db.League, error) {
	return s.db.GetLeagues(ctx)
}
