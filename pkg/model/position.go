package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Position struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}
