package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseID struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
}

func NullResponseID() ResponseID {
	return ResponseID{Id: primitive.NilObjectID}
}

func NewResponseID(id_ primitive.ObjectID) ResponseID {
	return ResponseID{Id: id_}
}

type Response struct {
	Status  string `json:"status" bson:"status"`
	Message string `json:"error" bson:"error"`
	Data    any    `json:"data,omitempty" bson:"data,omitempty"`
}

func NewSuccessResponse(status_ string, msg_ string, data_ any) Response {
	return Response{Status: status_, Message: msg_, Data: data_}
}

func NewErrorResponse(status_ string, msg_ string) Response {
	return Response{Status: status_, Message: msg_}
}

type ResponseToken struct {
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	AccessToken  string `json:"access_token" bson:"access_token"`
}
