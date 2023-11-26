package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Location string             `json:"location" bson:"location"`
	Name     string             `json:"name" bson:"name"`
}
