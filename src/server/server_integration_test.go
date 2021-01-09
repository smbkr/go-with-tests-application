package server

import (
	"application/src/data"
	"application/src/model"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingAndRetreivingWins(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mongoStore, disconnect := data.NewMongoPlayerStore(context.TODO())
	defer disconnect(context.TODO())

	var dataStores = map[string]PlayerStore{
		"in memory": data.NewInMemoryPlayerStore(),
		"mongodb":   mongoStore,
	}

	player := "Pepper"

	for name, store := range dataStores {
		server := NewPlayerServer(store)

		server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))
		server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))
		server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))

		t.Run(fmt.Sprintf("%s get score", name), func(t *testing.T) {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, playerScoreRequest(player))
			assertStatus(t, response, http.StatusOK)
			assertScoreResponse(t, response, 3)
		})

		t.Run(fmt.Sprintf("%s get league", name), func(t *testing.T) {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, leagueRequest())
			assertStatus(t, response, http.StatusOK)

			got := decodeLeagueResponse(t, response)
			want := []model.Player{
				{"Pepper", 3},
			}
			assertLeagueMatches(t, got, want)
		})
	}
}
