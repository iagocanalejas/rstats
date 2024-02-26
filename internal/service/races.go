package service

import (
	"log"

	"github.com/iagocanalejas/regatas/internal/service/races"
)

func (s *Service) SearchRaces(keywords string) []races.Race {
	log.Println("searching for races with: ", keywords)

	flatRaces, err := s.db.SearchRaces(keywords)
	if err != nil {
		log.Println("error loading races: ", err)
		return make([]races.Race, 0)
	}

	rs := make([]races.Race, len(flatRaces))
	for idx, race := range flatRaces {
		rs[idx] = *races.New(race)
	}
	return rs
}
