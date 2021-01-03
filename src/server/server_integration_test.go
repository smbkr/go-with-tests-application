package server

import (
	"application/src/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingAndRetrievingWins(t *testing.T) {
	store := data_store.NewInMemoryPlayerStore()
	server := PlayerServer{store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), recordWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, playerScoreRequest(player))
	assertStatus(t, response, http.StatusOK)
	assertScoreResponse(t, response, 3)
}
