package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Country string             `json:"country" bson:"country"`
	City    string             `json:"city" bson:"city"`
	Name    string             `json:"name" bson:"name"`
	Rooms   []Room             `json:"rooms" bson:"rooms"`
}

type HotelPostParams struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

type Room struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Price   float64            `json:"price" bson:"price"`
	SeaView bool               `json:"seaView" bson:"seaView"`
	BedType string             `json:"bedType" bson:"bedType"`
	Premium string             `json:"premium" bson:"premium"`
	HotelId primitive.ObjectID `json:"hotelId" bson:"hotelId"`
}
type RoomPostParams struct {
	Price   float64            `json:"price"`
	SeaView bool               `json:"seaView"`
	BedType string             `json:"bedType"`
	Premium string             `json:"premium"`
	HotelId primitive.ObjectID `json:"hotelId"`
}
