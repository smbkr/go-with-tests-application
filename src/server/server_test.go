package server

import (
	"application/src/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetPlayerScore(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := NewPlayerServer(store)

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

func TestAddPlayerScore(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{},
	}
	server := NewPlayerServer(store)

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

func TestLeague(t *testing.T) {
	store := &StubPlayerStore{}
	server := NewPlayerServer(store)

	t.Run("it returns a list of players", func(t *testing.T) {
		expectedLeague := []model.Player{
			{"Cleo", 30},
			{"Chris", 20},
			{"Charlie", 10},
		}
		store.league = expectedLeague
		response := httptest.NewRecorder()
		request, _ := leagueRequest()

		server.ServeHTTP(response, request)

		got := decodeLeagueResponse(t, response)
		assertStatus(t, response, http.StatusOK)
		assertJsonHeader(t, response, "application/json")
		assertLeagueMatches(t, got, expectedLeague)
	})
}

func assertJsonHeader(t *testing.T, response *httptest.ResponseRecorder, expected string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != expected {
		t.Errorf("response has incorrect content type, got %q want %q", response.Result().Header, expected)
	}
}

func assertLeagueMatches(t *testing.T, got []model.Player, expectedLeague []model.Player) {
	t.Helper()
	if !reflect.DeepEqual(got, expectedLeague) {
		t.Errorf("got %v, want %v", got, expectedLeague)
	}
}

func decodeLeagueResponse(t *testing.T, response *httptest.ResponseRecorder) []model.Player {
	t.Helper()
	var got []model.Player
	err := json.NewDecoder(response.Body).Decode(&got)
	if err != nil {
		t.Fatalf("unable to parse response %q, %v", response.Body, err)
	}
	return got
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

func leagueRequest() (*http.Request, error) {
	return http.NewRequest(http.MethodGet, "/league", nil)
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []model.Player
}

func (s *StubPlayerStore) PlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) League() []model.Player {
	return s.league
}
