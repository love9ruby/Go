package main

import (
	"db/model"
	"db/service/user"
	"fmt"

	"github.com/google/uuid"
)

func processUser(userService *user.UserService) {

	// init Tables
	if err := (*userService).InitTable(); err != nil {
		panic(err)
		return
	}

	// Create
	for i := 0; i < 10; i++ {
		randomStr := uuid.New().String()
		newName := fmt.Sprintf("test%s", randomStr)
		newEmail := fmt.Sprintf("test%s@abc.com", randomStr)
		userTmp := &model.User{
			Name:  newName,
			Email: newEmail,
		}
		if err := (*userService).CreateUser(userTmp.Email, userTmp.Name); err != nil {
			panic(err)
		}
	}

	// Read
	users, err := (*userService).GetUserList()
	if err != nil {
		panic(err)
		return
	}
	{
		for _, v := range users {
			println(v.Name, v.Email)
		}
	}

	// Update
	users[0].Email = "newEmail-" + uuid.NewString()
	if err := (*userService).UpdateUser(users[0].ID, users[0].Email); err != nil {
		panic(err)
		return
	}

	// Read
	users, err = (*userService).GetUserList()
	if err != nil {
		panic(err)
		return
	}
	{
		for _, v := range users {
			println(v.Name, v.Email)
		}
	}

	// Delete
	{
		for _, v := range users {
			if err := (*userService).DeleteUser(v.ID); err != nil {
				panic(err)
				return
			}
		}
	}

	// Read
	println("After delete")
	users, err = (*userService).GetUserList()
	if err != nil {
		panic(err)
		return
	}
	{
		for _, v := range users {
			println(v.Name)
		}
	}

	// Close
	user.Close(*userService)
}

func main() {
	var mode user.UserServiceType
	mode = user.Rdb
	//mode = user.Postgress
	//mode = user.PgOrm
	//mode = user.Mongo
	userService := user.NewUserService(mode)

	processUser(userService)
}
