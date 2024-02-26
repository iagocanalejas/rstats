package service

import (
	"context"
	"time"

	"github.com/iagocanalejas/regatas/internal/db"
)

type Service struct {
	db db.Repository
}

func Init() *Service {
	return &Service{
		db: db.New(),
	}
}

func (s *Service) IsHealthy() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return s.db.IsHealthy(ctx)
}
