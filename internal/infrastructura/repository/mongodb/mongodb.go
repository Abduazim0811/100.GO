package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"100.GO/internal/entity/origin"
	"100.GO/internal/entity/user"
	"100.GO/internal/infrastructura/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongodb struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewUserMongodb(client *mongo.Client, collection *mongo.Collection) repository.UserRepository {
	return &UserMongodb{client: client, collection: collection}
}

func (u *UserMongodb) AddUser(user *user.CreateUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := u.collection.FindOne(ctx, bson.M{"email": user.Email})
	if res.Err() == nil {
		return fmt.Errorf("email already registered")
	}
	if res.Err() != mongo.ErrNoDocuments {
		return res.Err()
	}

	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error adding user:", err)
		return err
	}
	return nil
}

func (u *UserMongodb) GetUserByEmail(email string) (*user.Login, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := u.collection.FindOne(ctx, bson.M{"email": email})

	var user user.Login
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserMongodb) AddOrigin(req origin.CreateOrigin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := u.collection.InsertOne(ctx, req)
	if err != nil {
		log.Println("Error adding origin:", err)
		return err
	}
	return nil
}

func (u *UserMongodb) GetByIdOrigin(reqId string) (*origin.GetOrigin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(reqId)
	if err != nil {
		return nil, fmt.Errorf("invalid id format")
	}

	res := u.collection.FindOne(ctx, bson.M{"_id": id})

	var origin origin.GetOrigin
	err = res.Decode(&origin)
	if err != nil {
		return nil, err
	}

	return &origin, nil
}

func (u *UserMongodb) UpdateOrigin(id primitive.ObjectID, req origin.CreateOrigin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"origin": req.Origin}}
	var ord origin.GetOrigin
	err := u.collection.FindOneAndUpdate(ctx, filter, update).Decode(&ord)
	if err != nil {
		log.Println("Error updating origin:", err)
		return err
	}
	return nil
}

func (u *UserMongodb) DeleteOrigin(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := u.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Println("Error deleting origin:", err)
		return err
	}
	return nil
}

func (u *UserMongodb) GetAllOrigins() ([]*origin.GetOrigin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := u.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error fetching origins:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var origins []*origin.GetOrigin
	for cursor.Next(ctx) {
		var origin origin.GetOrigin
		if err := cursor.Decode(&origin); err != nil {
			log.Println("Error decoding origin:", err)
			return nil, err
		}
		origins = append(origins, &origin)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return origins, nil
}

func (u *UserMongodb) OriginGetAll() ([]*origin.CreateOrigin, error) {
	val, err := u.GetAllOrigins()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var origins []*origin.CreateOrigin

	for _, v := range val {
		var all origin.CreateOrigin
		data, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if err := json.Unmarshal(data, &all); err != nil {
			log.Println(err)
			return nil, err
		}
		origins = append(origins, &all)
	}
	return origins, nil
}
