package user

import (
	"db/model"
	"db/repository/orm"
	"gorm.io/gorm"
)

type PgOrmUserService struct {
	alias UserServiceType
	gorm  *gorm.DB
}

func (pgOrm *PgOrmUserService) InitTable() error {
	err := orm.CreateTable()
	if err != nil {
		return err
	}
	return nil
}

func (pgOrm *PgOrmUserService) GetType() UserServiceType {
	return pgOrm.alias
}

func (pgOrm *PgOrmUserService) CreateUser(email, name string) error {
	err := orm.InsertRecord(name, email)
	if err != nil {
		return err
	}
	return nil
}

func (pgOrm *PgOrmUserService) UpdateUser(id, email string) error {
	err := orm.UpdateRecord(id, email)
	if err != nil {
		return err
	}
	return nil
}

func (pgOrm *PgOrmUserService) DeleteUser(id string) error {
	err := orm.DeleteRecord(id)
	return err
}

func (pgOrm *PgOrmUserService) GetUserList() ([]*model.User, error) {
	users, err := orm.QueryRecords()
	if err != nil {
		return nil, err
	}
	return users, nil
}
