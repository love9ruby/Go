package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	MongoClient *mongo.Client
	DB          *mongo.Database
	ctx         context.Context
}

type Url struct {
	ID       string `bson:"_id,omitempty"`
	ShortURL string `bson:"short_url"`
	LongURL  string `bson:"long_url"`
	Password string `bson:"password"`
}

type key struct {
	ID  string `bson:"_id,omitempty"`
	Key string `bson:"key"`
}
