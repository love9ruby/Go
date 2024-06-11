package user

import (
	"database/sql"
	"db/model"
	"db/repository/pg"
	"gorm.io/gorm"
)

type PostgressUserService struct {
	db    *sql.DB
	gorm  *gorm.DB
	alias UserServiceType
}

func (pgs *PostgressUserService) InitTable() error {
	err := pg.CreateTable(pgs.db)
	if err != nil {
		return err
	}
	return nil
}

func (pgs *PostgressUserService) GetType() UserServiceType {
	return pgs.alias
}

func (pgs *PostgressUserService) CreateUser(email, name string) error {
	err := pg.InsertRecord(pgs.db, name, email)
	if err != nil {
		return err
	}
	return nil
}

func (pgs *PostgressUserService) UpdateUser(id, email string) error {
	err := pg.UpdateRecord(pgs.db, id, email)
	if err != nil {
		return err
	}
	return nil
}

func (pgs *PostgressUserService) DeleteUser(id string) error {
	err := pg.DeleteRecord(pgs.db, id)
	if err != nil {
		return err
	}
	return nil
}

func (pgs *PostgressUserService) GetUserList() ([]*model.User, error) {
	users, err := pg.QueryRecords(pgs.db)
	if err != nil {
		return nil, err
	}
	return users, nil
}
