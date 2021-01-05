package server

import (
	"application/src/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := &PlayerServer{}
	p.store = store
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playerHandler))
	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(p.store.League())
}

func (p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.playerWinHandler(w, player)
	case http.MethodGet:
		p.playerScoreHandler(w, player)
	}
}

func (p *PlayerServer) playerWinHandler(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) playerScoreHandler(w http.ResponseWriter, player string) {
	score := p.store.PlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprint(w, score)
}

type PlayerStore interface {
	PlayerScore(name string) int
	RecordWin(name string)
	League() []model.Player
}
