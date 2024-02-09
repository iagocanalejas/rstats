package service

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

func (s *Service) GetFlags(ctx context.Context) ([]db.Flag, error) {
	return s.db.GetFlags(ctx)
}

func (s *Service) GetTrophies(ctx context.Context) ([]db.Trophy, error) {
	return s.db.GetTrophies(ctx)
}
