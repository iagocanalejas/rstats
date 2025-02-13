package service

import (
	"github.com/iagocanalejas/rstats/internal/db"
	"github.com/iagocanalejas/rstats/internal/types"
	prettylog "github.com/iagocanalejas/rstats/internal/utils/pretty-log"
)

type Service struct {
	db db.Repository
}

func Init() *Service {
	return &Service{
		db: db.New(),
	}
}

func Static() *Service {
	// This is a static service that doesn't use the database.
	return &Service{}
}

func (s *Service) GetLeagueByID(leagueID int64) (*types.League, error) {
	dbLeague, err := s.db.GetLeagueByID(leagueID)
	if err != nil {
		prettylog.Error("error loading league: %v", err)
		return nil, err
	}

	l := types.NewLeagueFromDB(dbLeague)
	return l, nil
}

func (s *Service) GetFlagByID(flagID int64) (*types.Flag, error) {
	dbFlag, err := s.db.GetFlagByID(flagID)
	if err != nil {
		prettylog.Error("error loading flag: %v", err)
		return nil, err
	}

	f := types.NewFlagFromDB(dbFlag, nil)
	return f, nil
}

func (s *Service) GetClubByID(clubID int64) (*types.Entity, error) {
	dbClub, err := s.db.GetClubByID(clubID)
	if err != nil {
		prettylog.Error("error loading club: %v", err)
		return nil, err
	}

	e := types.NewEntityFromDB(dbClub, nil)
	return e, nil
}
