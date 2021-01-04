package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	server := &PlayerServer{
		Store: &StubPlayerStore{
			scores: map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
		},
	}

	t.Run(`returns score for "Pepper"`, func(t *testing.T) {
		response := httptest.NewRecorder()
		request := playerScoreRequest("Pepper")
		server.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
		assertScoreResponse(t, response, 20)
	})

	t.Run(`returns the score for "Floyd"`, func(t *testing.T) {
		response := httptest.NewRecorder()
		request := playerScoreRequest("Floyd")
		server.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)
		assertScoreResponse(t, response, 10)
	})

	t.Run("returns 404 for a non-existent player", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := playerScoreRequest("Foobar")
		server.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestAddScore(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{},
	}
	server := &PlayerServer{&store}

	t.Run("records wins", func(t *testing.T) {
		player := "Pepper"
		response := httptest.NewRecorder()
		request := recordWinRequest(player)

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("RecordWin called %v times, expected 1", len(store.winCalls))
		}

		if store.winCalls[0] != player {
			t.Errorf("RecordWin called with incorrect player name %q, expected %q", store.winCalls[0], player)
		}
	})
}

func assertStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("incorrect status code: got %v want %v", got, want)
	}
}

func assertScoreResponse(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Body.String()
	if got != fmt.Sprint(want) {
		t.Errorf("unexpected response: got %s want %d", got, want)
	}
}

func playerScoreRequest(p string) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", p), nil)
	return r
}

func recordWinRequest(p string) *http.Request {
	r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", p), nil)
	return r
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) (int, error) {
	score, found := s.scores[name]
	if !found {
		return 0, errors.New("not found")
	}
	return score, nil
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
