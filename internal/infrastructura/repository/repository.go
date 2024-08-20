package repository

import (
	"100.GO/internal/entity/origin"
	"100.GO/internal/entity/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type UserRepository interface{
	AddUser(user *user.CreateUser)error
	GetUserByEmail(email string)(*user.Login, error)
	AddOrigin(req origin.CreateOrigin)error
	GetByIdOrigin(reqId string)(*origin.GetOrigin, error)
	GetAllOrigins()([]*origin.GetOrigin, error)
	UpdateOrigin(id primitive.ObjectID, req origin.CreateOrigin)error
	DeleteOrigin(id primitive.ObjectID) error
	OriginGetAll()([]*origin.CreateOrigin, error)
}