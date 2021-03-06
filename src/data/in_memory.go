package data

import "application/src/model"

type InMemoryPlayerStore struct {
	scores map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		make(map[string]int),
	}
}

func (s *InMemoryPlayerStore) PlayerScore(name string) int {
	return s.scores[name]
}

func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.scores[name]++
}

func (s *InMemoryPlayerStore) League() []model.Player {
	var league []model.Player
	for name, wins := range s.scores {
		league = append(league, model.Player{Name: name, Wins: wins})
	}
	return league
}
