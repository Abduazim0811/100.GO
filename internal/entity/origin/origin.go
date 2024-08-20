package origin

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateOrigin struct {
	Origin string `json:"origin" bson:"origin"`
}

type GetOrigin struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Origin string 		  `json:"origin" bson:"origin"`
}
