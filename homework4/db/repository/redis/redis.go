package redis

import (
	"context"
	"db/model"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client

func GetDB() *redis.Client {
	if rdb != nil {
		return rdb
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return rdb
}

func Close() {
	err := rdb.Close()
	if err != nil {
		panic(err)
	}
	rdb = nil
}

func InsertRecord(name, email string) error {
	key := uuid.NewString()
	value := &model.RedisUser{
		Name:  name,
		Email: email,
	}
	// Marshal the user to jsonStr
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = rdb.Set(context.Background(), key, string(jsonValue), time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func UpdateRecord(id, email string) error {
	value, err := rdb.Get(context.Background(), id).Result()
	if err != nil {
		return err
	}
	user := &model.RedisUser{}
	// Unmarshal the value to user
	err = json.Unmarshal([]byte(value), user)
	user.Email = email
	// Marshal the user to jsonStr
	jsonValue, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = rdb.Set(context.Background(), id, string(jsonValue), time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func DeleteRecord(id string) error {
	err := rdb.Del(context.Background(), id).Err()
	if err != nil {
		return err
	}
	return nil
}

func QueryRecords() ([]*model.User, error) {
	keys, err := rdb.Keys(context.Background(), "*").Result()
	if err != nil {
		return nil, err
	}
	users := make([]*model.User, 0)
	for _, key := range keys {
		value, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		originUser := &model.RedisUser{}
		// Unmarshal the value to user
		err = json.Unmarshal([]byte(value), originUser)
		if err != nil {
			return nil, err
		}
		user := &model.User{
			ID:    key,
			Name:  originUser.Name,
			Email: originUser.Email,
		}
		users = append(users, user)
	}
	return users, nil
}
