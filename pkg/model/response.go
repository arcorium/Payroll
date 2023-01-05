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

type ErrorResponse struct {
	Status  string `json:"status" bson:"status"`
	Message string `json:"error" bson:"error"`
}

type SuccessResponse struct {
	Status string `json:"status" bson:"status"`
	Data   any    `json:"data" bson:"data"`
}

func NewSuccessResponse(status_ string, data_ any) SuccessResponse {
	return SuccessResponse{Status: status_, Data: data_}
}

func NewErrorResponse(status_ string, msg_ string) ErrorResponse {
	return ErrorResponse{Status: status_, Message: msg_}
}

type ResponseToken struct {
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	AccessToken  string `json:"access_token" bson:"access_token"`
}
