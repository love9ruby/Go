package orm

import (
	"db/model"
	"fmt"
	"github.com/google/uuid"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

var gormDB *gorm.DB

func GetDB() *gorm.DB {
	if gormDB != nil {
		return gormDB
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gormDB = db
	return gormDB
}

func Close() {
	db, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
	gormDB = nil
}

func CreateTable() error {
	err := gormDB.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}
	return nil
}

func InsertRecord(name, email string) error {
	user := model.User{Name: name, Email: email, ID: uuid.New().String()}
	result := gormDB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateRecord(id, email string) error {
	err := gormDB.Model(&model.User{}).Where("id = ?", id).Update("email", email).Error
	return err
}

func DeleteRecord(id string) error {
	err := gormDB.Where("id = ?", id).Delete(&model.User{}).Error
	return err
}

func QueryRecords() ([]*model.User, error) {
	users := make([]*model.User, 0)
	result := gormDB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
