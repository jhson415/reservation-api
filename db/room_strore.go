package db

import (
	"context"
	"log"

	"github.com/jhson415/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	Dropper
	PostRoom(context.Context, *types.Room) (*types.Room, error)
	GetRoomListByHotelId(context.Context, string) (*[]types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	coll       *mongo.Collection
	HotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DB_NAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}
func (m MongoRoomStore) PostRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := m.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.Id = res.InsertedID.(primitive.ObjectID)
	return room, nil
}

func (m MongoRoomStore) GetRoomListByHotelId(ctx context.Context, id string) (*[]types.Room, error) {
	var roomList []types.Room
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatalln("Error while converting id to objectID", err)
		return nil, err
	}
	cur, err := m.coll.Find(ctx, bson.M{"hotelId": oId})
	if err != nil {
		log.Fatalln("Error while creating cursor", err)
		return nil, err
	}
	err = cur.All(ctx, &roomList)
	if err != nil {
		log.Fatalln("Failed to decode cursor", err)
	}
	return &roomList, nil

}

func (m MongoRoomStore) Drop(ctx context.Context) error {
	err := m.coll.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}
