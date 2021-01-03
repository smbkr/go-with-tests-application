package server

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerServer struct {
	Store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.handleWin(w, player)
	case http.MethodGet:
		p.handlePlayerScore(w, player)
	}
}

func (p *PlayerServer) handleWin(w http.ResponseWriter, player string) {
	p.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
	return
}

func (p *PlayerServer) handlePlayerScore(w http.ResponseWriter, player string) {
	score, err := p.Store.GetPlayerScore(player)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprint(w, score)
}

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	RecordWin(name string)
}
