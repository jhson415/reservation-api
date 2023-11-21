package db

import (
	"context"

	"github.com/jhson415/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {

	return &MongoUserStore{
		client: client,
		coll:   client.Database(DB_NAME).Collection(userColl),
	}
}

func (m MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user = types.User{}
	if err = m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
