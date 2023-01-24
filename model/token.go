package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Token  string             `json:"refresh_token" bson:"refresh_token"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
}
