package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        	primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
	Firstname 	string             	`json:"firstname" bson:"firstname"`
	Lastname	string				`json:"lastname" bson:"lastname"`
	Email		string				`json:"email" bson:"email"`
	Password    string				`json:"password" bson:"password"`
	Code 		int					`json:"code" bson:"code"`
}

type CreateUser struct{
	Firstname 	string             	`json:"firstname" bson:"firstname"`
	Lastname	string				`json:"lastname" bson:"lastname"`
	Email		string				`json:"email" bson:"email"`
	Password    string				`json:"password" bson:"password"`
}

type Login struct{
	Email		string				`json:"email" bson:"email"`
	Password    string				`json:"password" bson:"password"`
}

type VerifyCode struct{
	Email		string				`json:"email" bson:"email"`
	Code 		int					`json:"code" bson:"code"`
}