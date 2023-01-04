package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Teach struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StaffId   primitive.ObjectID `json:"staff_id" bson:"staff_id"`
	Institute string             `json:"institute" bson:"institute"`
	Details   []TeachDetail      `json:"details" bson:"details"`
}

type TeachDetail struct {
	Study string `json:"study" bson:"study"`
	Hours uint8  `json:"hours" bson:"hours"`
}
