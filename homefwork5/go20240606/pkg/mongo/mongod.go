package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewClient(url string, defaultDB ...string) (*Client, error) {
	if len(defaultDB) == 0 {
		defaultDB = append(defaultDB, "url")
	}
	var err error
	client := new(Client)
	client.ctx = context.TODO()
	client.MongoClient, err = mongo.Connect(client.ctx, options.Client().ApplyURI(url).SetConnectTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	client.DB = client.MongoClient.Database(defaultDB[0])
	return client, nil
}

func (c *Client) Ping() error {
	if c.MongoClient == nil || c.ctx == nil {
		return errors.New("MongoClient is nil")
	}
	return c.MongoClient.Ping(c.ctx, nil)
}

func (c *Client) Close() error {
	if c.MongoClient == nil || c.ctx == nil {
		return errors.New("MongoClient is nil")
	}
	return c.MongoClient.Disconnect(context.Background())
}

func (c *Client) CreateUrl(url Url) error {
	_, err := c.DB.Collection("urls").InsertOne(c.ctx, url)
	return err
}

func (c *Client) FindUrl(shortURL string) (Url, error) {
	var url Url
	filter := bson.D{{"short_url", bson.D{{"$eq", shortURL}}}}
	err := c.DB.Collection("urls").FindOne(c.ctx, filter).Decode(&url)
	if err != nil {
		return Url{}, err
	}
	return url, err
}

func (c *Client) InitKeyIndex() error {
	col := c.DB.Collection("keys")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "key", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := col.Indexes().CreateOne(c.ctx, indexModel)
	return err
}

func (c *Client) CreateKey(newKey string) error {
	_, err := c.DB.Collection("keys").InsertOne(c.ctx, key{
		Key: newKey,
	})
	return err
}

func (c *Client) IsKeyExist(key string) (bool, error) {
	count, err := c.DB.Collection("keys").CountDocuments(c.ctx, bson.M{"key": key})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
