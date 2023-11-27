package db

import "context"

const DB_NAME = "hotel-reservation"
const DB_NAME_TEST = "hotel-reservation-test"
const DB_URI = "mongodb://localhost:27017"

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

type Dropper interface {
	Drop(ctx context.Context) error
}
