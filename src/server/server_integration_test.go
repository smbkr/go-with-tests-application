package server

import (
	"application/src/data"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerDataStoreIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mongoClient := mongoClient(ctx)
	defer mongoClient.Disconnect(ctx)
	_ = mongoClient.Database("game").Collection("players").Drop(ctx)

	var dataStores = map[string]PlayerStore{
		"in memory": data.NewInMemoryPlayerStore(),
		"mongodb":   &data.MongoPlayerStore{mongoClient},
	}

	for name, store := range dataStores {
		t.Run(name, func(t *testing.T) {
			server := PlayerServer{store}
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

func mongoClient(ctx context.Context) *mongo.Client {
	connectionURI := "mongodb://root:root@localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
