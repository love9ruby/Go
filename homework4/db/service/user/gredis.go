package user

import (
	"db/model"
	"github.com/redis/go-redis/v9"
)
import rdb "db/repository/redis"

type RedisUserService struct {
	alias  UserServiceType
	client *redis.Client
}

func (rds *RedisUserService) InitTable() error {
	// No need to create table in Redis
	return nil
}

func (rds *RedisUserService) GetType() UserServiceType {
	return rds.alias
}

func (rds *RedisUserService) CreateUser(email, name string) error {
	err := rdb.InsertRecord(name, email)
	if err != nil {
		return err
	}
	return nil
}

func (rds *RedisUserService) UpdateUser(id, email string) error {
	err := rdb.UpdateRecord(id, email)
	if err != nil {
		return err
	}
	return nil
}

func (rds *RedisUserService) DeleteUser(id string) error {
	err := rdb.DeleteRecord(id)
	if err != nil {
		return err
	}
	return nil
}

func (rds *RedisUserService) GetUserList() ([]*model.User, error) {
	users, err := rdb.QueryRecords()
	if err != nil {
		return nil, err
	}
	return users, nil
}
