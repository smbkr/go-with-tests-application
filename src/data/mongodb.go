package data

import (
	"application/src/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Player struct {
	Name string `bson:"name"`
	Wins int    `bson:"wins"`
}

type MongoPlayerStore struct {
	client *mongo.Client
}

func (s *MongoPlayerStore) PlayerScore(name string) int {
	player := Player{}
	err := s.client.
		Database("game").
		Collection("players").
		FindOne(context.TODO(), bson.M{"name": name}).
		Decode(&player)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
	}

	return player.Wins
}

func (s *MongoPlayerStore) RecordWin(name string) {
	ctx := context.TODO()
	collection := s.client.Database("game").Collection("players")
	filter := bson.M{"name": name}
	update := bson.M{"$inc": bson.M{"wins": 1}}
	result := collection.FindOneAndUpdate(ctx, filter, update)
	err := result.Err()
	switch err {
	case nil:
		return
	case mongo.ErrNoDocuments:
		player := Player{
			Name: name,
			Wins: 1,
		}
		res, err := collection.InsertOne(ctx, &player)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("added new player with ID %s", res.InsertedID)
	default:
		log.Fatal(err)
	}
}

func (s *MongoPlayerStore) League() []model.Player {
	var league []model.Player
	ctx := context.TODO()
	collection := s.client.Database("game").Collection("players")
	projection := bson.M{"name": true, "wins": true}
	opts := options.Find().SetProjection(projection)
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &league); err != nil {
		log.Fatal(err)
	}
	return league
}

func NewMongoPlayerStore(ctx context.Context) (store *MongoPlayerStore, disconnect func(context.Context)) {
	connectionURI := "mongodb://root:root@localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	disconnect = func(ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &MongoPlayerStore{client}, disconnect
}
