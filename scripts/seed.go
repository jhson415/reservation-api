package main

import (
	"context"
	"log"

	"github.com/jhson415/reservation-api/db"
	"github.com/jhson415/reservation-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	hotelList = []types.Hotel{
		{
			Country: "Korea",
			City:    "Seoul",
			Name:    "H for Living Hotel",
		},
		{
			Country: "Korea",
			City:    "California",
			Name:    "California Hotel",
		},
		{
			Country: "England",
			City:    "London",
			Name:    "big ben hotel",
		}}
	basicRoomList = []types.Room{
		{
			Price:   100.0,
			SeaView: false,
			BedType: "Single",
			Premium: "Basic",
		},
		{
			Price:   130.0,
			SeaView: true,
			BedType: "Double",
			Premium: "Basic",
		},
		{
			Price:   330.0,
			SeaView: true,
			BedType: "King",
			Premium: "Suite",
		},
	}
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func main() {
	seedHotel()

}

func seedHotel() {
	for _, hotel := range hotelList {
		_, err := hotelStore.PostHotel(context.TODO(), &hotel)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(hotel)
		//TODO Add room with hotel id
		for _, room := range basicRoomList {
			room.HotelId = hotel.Id
			_, err := roomStore.PostRoom(context.TODO(), &room)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(room)
		}
	}
}

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
