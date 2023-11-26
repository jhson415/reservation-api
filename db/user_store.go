package db

import (
	"context"
	"fmt"

	"github.com/jhson415/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStore interface {
	Dropper
	GetUserById(context.Context, string) (*types.User, error)
	GetUserList(context.Context) (*[]types.User, error)
	PostUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	PutUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
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

func (m MongoUserStore) GetUserList(ctx context.Context) (*[]types.User, error) {
	var (
		users = []types.User{}
	)
	m.coll.Find(ctx, bson.M{})
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return &users, nil
}

func (m MongoUserStore) PostUser(ctx context.Context, user *types.User) (*types.User, error) {

	result, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := m.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("Given User ID not found")
	}
	return nil
}

func (m MongoUserStore) PutUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.M{"$set": params.ToBson()}
	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil

}

func (m MongoUserStore) Drop(ctx context.Context) error {
	if err := m.coll.Drop(ctx); err != nil {
		return err
	}
	return nil
}
