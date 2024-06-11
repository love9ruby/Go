package user

import (
	"db/model"
	"go.mongodb.org/mongo-driver/mongo"
)
import mongoRepository "db/repository/mongo"

type MongoUserService struct {
	Client *mongo.Client
	DB     *mongo.Database
	alias  UserServiceType
}

func (mgs *MongoUserService) InitTable() error {
	err := mongoRepository.CreateTable(mgs.DB)
	if err != nil {
		return err
	}
	return nil
}

func (mgs *MongoUserService) GetType() UserServiceType {
	return mgs.alias
}

func (mgs *MongoUserService) CreateUser(email, name string) error {
	err := mongoRepository.InsertRecord(mgs.DB, name, email)
	if err != nil {
		return err
	}
	return nil
}

func (mgs *MongoUserService) UpdateUser(id, email string) error {
	err := mongoRepository.UpdateRecord(mgs.DB, id, email)
	if err != nil {
		return err
	}
	return nil
}

func (mgs *MongoUserService) DeleteUser(id string) error {
	err := mongoRepository.DeleteRecord(mgs.DB, id)
	if err != nil {
		return err
	}
	return nil
}

func (mgs *MongoUserService) GetUserList() ([]*model.User, error) {
	users, err := mongoRepository.QueryRecords(mgs.DB)
	if err != nil {
		return nil, err
	}
	return users, nil
}
