package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Config struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SecretKey string             `json:"secret_key" bson:"secret_key"`
}
