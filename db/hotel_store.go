package db

import (
	"context"
	"log"

	"github.com/jhson415/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	Dropper
	PostHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	GetHotelList(context.Context) (*[]types.Hotel, error)
	GetHotelById(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DB_NAME).Collection(hotelColl),
	}
}

func (m MongoHotelStore) PostHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := m.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (m MongoHotelStore) GetHotelList(ctx context.Context) (*[]types.Hotel, error) {
	var hotels []types.Hotel
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return &hotels, nil
}

func (m MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	var hotel types.Hotel
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatalln("Error on converting string to oid", err)
		return nil, err
	}
	if err = m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		log.Fatalln("Error on Finding all the documents", err)
		return nil, err
	}
	return &hotel, nil
}

func (m MongoHotelStore) Drop(ctx context.Context) error {
	err := m.coll.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}
