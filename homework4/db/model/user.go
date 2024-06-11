package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

type MongoUser struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
}

type RedisUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
