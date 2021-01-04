package server

import (
	"application/src/data"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerDataStoreIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mongoStore, disconnect := data.NewMongoPlayerStore(context.TODO())
	defer disconnect(context.TODO())

	var dataStores = map[string]PlayerStore{
		"in memory": data.NewInMemoryPlayerStore(),
		"mongodb":   mongoStore,
	}

	for name, store := range dataStores {
		t.Run(name, func(t *testing.T) {
			server := NewPlayerServer(store)
			player := "Pepper"

			server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))
			server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))
			server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))

			response := httptest.NewRecorder()
			server.ServeHTTP(response, playerScoreRequest(player))
			assertStatus(t, response, http.StatusOK)
			assertScoreResponse(t, response, 3)
		})
	}
}
