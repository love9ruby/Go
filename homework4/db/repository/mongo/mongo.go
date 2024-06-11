package mongo

import (
	"context"
	"db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var MongoClient *mongo.Client

func GetDB() (*mongo.Client, *mongo.Database) {
	if MongoClient != nil {
		return MongoClient, MongoClient.Database("test")
	}
	mongoUri := "mongodb://localhost:27017"
	var err error
	MongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri).SetConnectTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}
	db := MongoClient.Database("test")
	return MongoClient, db
}

func Close() {
	if err := MongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	MongoClient = nil
}

func CreateTable(db *mongo.Database) error {
	// mongoDB does not require to create table
	return nil
}

func InsertRecord(db *mongo.Database, name, email string) error {
	// mongoDB insert one user
	collection := db.Collection("users")
	_, err := collection.InsertOne(context.TODO(), model.MongoUser{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateRecord(db *mongo.Database, id string, email string) error {
	// mongoDB update one user
	collection := db.Collection("users")
	// string to primitive.ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{{"email", email}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRecord(db *mongo.Database, id string) error {
	// mongoDB delete one user
	collection := db.Collection("users")
	objectId, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objectId}}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func QueryRecords(db *mongo.Database) ([]*model.User, error) {
	// mongoDB query all users
	collection := db.Collection("users")
	cursor, err := collection.Find(context.TODO(), primitive.D{{}})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor) {
		err := cursor.Close(context.Background())
		if err != nil {
			panic(err)
		}
	}(cursor)

	users := make([]*model.User, 0)
	for cursor.Next(context.Background()) {
		originUser := model.MongoUser{}
		err := cursor.Decode(&originUser)
		if err != nil {
			return nil, err
		}
		user := model.User{
			ID:    originUser.ID.Hex(),
			Name:  originUser.Name,
			Email: originUser.Email,
		}
		users = append(users, &user)
	}
	return users, nil
}
