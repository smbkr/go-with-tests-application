package data_store

import "errors"

type InMemoryPlayerStore struct {
	scores map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		make(map[string]int),
	}
}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) (int, error) {
	score, found := s.scores[name]
	if !found {
		return 0, errors.New("not found")
	}
	return score, nil
}

func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.scores[name]++
}
