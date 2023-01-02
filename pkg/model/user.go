package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type        UserType           `json:"type" bson:"type"`
	TeacherName string             `json:"teacher_name" bson:"teacher_name"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password" bson:"password"`
	IsLoggedIn  bool               `json:"is_logged_in" bson:"is_logged_in"`
	TeacherId   primitive.ObjectID `json:"detail" bson:"detail"`
}

func (u *User) SetDefaultValue(teacherId_ primitive.ObjectID) {
	u.Type = Admin
	u.IsLoggedIn = false
	u.TeacherId = teacherId_

	if len(u.Username) < 1 {
		u.Username = u.TeacherName
	}
}

type UserType string

const (
	Admin UserType = "admin"
	Super UserType = "super"
)
