package pg

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Client struct {
	Db *gorm.DB
}

type Dsn struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

func NewClient(dsn Dsn) (*Client, error) {
	// create dsn string
	dsnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dsn.Host, dsn.Port, dsn.User, dsn.Password, dsn.DB)
	db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Client{Db: db}, nil
}

func (c *Client) Close() error {
	db, err := c.Db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (c *Client) Migration() error {
	err := c.Db.AutoMigrate(&Event{}, &User{})
	if err != nil {
		return err
	}
	// create constraint for foreign key, events.password -> users.password
	var count int64
	err = c.Db.Raw("SELECT count(*) FROM information_schema.table_constraints WHERE constraint_name = 'fk_password'").Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	err = c.Db.Exec("ALTER TABLE events ADD CONSTRAINT fk_password FOREIGN KEY (password) REFERENCES users(password)").Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateEvent(event Event) error {
	tx := c.Db.Create(&event)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (c *Client) CreateUser(email string) (string, error) {
	uid := uuid.NewString()
	if tx := c.Db.Create(&User{
		Email:    email,
		Password: uid,
	}); tx.Error != nil {
		return "", tx.Error
	}
	return uid, nil
}

func (c *Client) FindEventByUser(password string) ([]Event, error) {
	var events []Event
	// query for events where password = password
	tx := c.Db.Where("password = ?", password).Find(&events)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return events, nil
}
