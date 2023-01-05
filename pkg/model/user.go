package model

import (
	"Penggajian/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type       UserType           `json:"type,omitempty" bson:"type"`
	Username   string             `json:"username,omitempty" bson:"username"`
	Password   string             `json:"password" bson:"password"`
	IsLoggedIn bool               `json:"is_logged_in,omitempty" bson:"is_logged_in"`

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
