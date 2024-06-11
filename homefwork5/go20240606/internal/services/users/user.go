package users

import (
	"go20240606/pkg/pg"
	"os"
)

var dsn = pg.Dsn{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "mysecretpassword",
	DB:       "postgres",
}

type IUserService interface {
	// CreateUser create a new user
	CreateUser(email string) (string, error)
	// GetUserStatistics User Statistics
	GetUserStatistics(email string) ([]pg.Event, error)
	// CreateEvent create a new event for a user
	CreateEvent(event pg.Event) error
}

type UserService struct {
	db *UserPgImpl
}

type UserPgImpl struct {
	db *pg.Client
}

func (up *UserPgImpl) CreateUser(email string) (string, error) {
	pwd, err := up.db.CreateUser(email)
	if err != nil {
		return "", err
	}
	return pwd, nil
}

func (up *UserPgImpl) GetUserStatistics(pwd string) ([]pg.Event, error) {
	events, err := up.db.FindEventByUser(pwd)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (up *UserPgImpl) CreateEvent(event pg.Event) error {
	err := up.db.CreateEvent(event)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService() *UserService {
	host := os.Getenv("POSTGRES_HOST")
	if host != "" {
		dsn.Host = host
	}
	pwd := os.Getenv("POSTGRES_PASSWORD")
	if pwd != "" {
		dsn.Password = pwd
	}
	db, err := pg.NewClient(dsn)
	if err != nil {
		panic(err)
	}
	err = db.Migration()
	if err != nil {
		panic(err)
	}
	return &UserService{
		db: &UserPgImpl{
			db: db,
		},
	}
}

func (us *UserService) CreateUser(email string) (string, error) {
	return us.db.CreateUser(email)
}

func (us *UserService) GetUserStatistics(pwd string) ([]pg.Event, error) {
	return us.db.GetUserStatistics(pwd)
}

func (us *UserService) CreateEvent(event pg.Event) error {
	return us.db.CreateEvent(event)
}
