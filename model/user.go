package model

import (
	"Penggajian/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type       UserType           `json:"type,omitempty" bson:"type,omitempty"`
	Username   string             `json:"username,omitempty" bson:"username,omitempty"`
	Password   string             `json:"password" bson:"password,omitempty"`
	IsLoggedIn bool               `json:"is_logged_in,omitempty" bson:"is_logged_in,omitempty"`

	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}

func (u *User) SetDefaultValue() {
	if util.IsEmpty(string(u.Type)) {
		u.Type = Admin
	}
	u.IsLoggedIn = false
	u.CreatedAt = time.Now()
	u.UpdateModifiedTime()
}

func (u *User) UpdateModifiedTime() {
	u.ModifiedAt = time.Now()
}

type UserType string

const (
	Admin UserType = "admin"
	Super UserType = "super"
)
