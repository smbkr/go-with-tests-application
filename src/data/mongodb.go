package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Player struct {
	Name  string `bson:"name"`
	Score int    `bson:"score"`
}

type MongoPlayerStore struct {
	Client *mongo.Client
}

func (s *MongoPlayerStore) GetPlayerScore(name string) int {
	player := Player{}
	err := s.Client.
		Database("game").
		Collection("players").
		FindOne(context.TODO(), bson.M{"name": name}).
		Decode(&player)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
	}

	return player.Score
}

func (s *MongoPlayerStore) RecordWin(name string) {
	ctx := context.TODO()
	collection := s.Client.
		Database("game").
		Collection("players")
	filter := bson.M{"name": name}
	update := bson.M{"$inc": bson.M{"score": 1}}
	result := collection.FindOneAndUpdate(ctx, filter, update)
	err := result.Err()
	switch err {
	case nil:
		return
	case mongo.ErrNoDocuments:
		player := Player{
			Name:  name,
			Score: 1,
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
