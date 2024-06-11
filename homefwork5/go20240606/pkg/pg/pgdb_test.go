package pg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var rightDsn = Dsn{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "mysecretpassword",
	DB:       "postgres",
}

func tearDown() {
	client, _ := NewClient(rightDsn)
	client.Db.Exec("DROP TABLE events")
	client.Db.Exec("DROP TABLE users")
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(rightDsn)
	assert.Nil(t, err)

	err = client.Close()
	assert.Nil(t, err)

	client, err = NewClient(Dsn{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "wrongpassword",
		DB:       "test",
	})
	assert.NotNil(t, err)
}

func TestClient_Migration(t *testing.T) {
	defer tearDown()
	client, err := NewClient(rightDsn)
	assert.Nil(t, err)

	err = client.Migration()
	assert.Nil(t, err)

	err = client.Close()
	assert.Nil(t, err)
}

func TestClient_CreateEvent(t *testing.T) {
	defer tearDown()
	client, err := NewClient(rightDsn)
	assert.Nil(t, err)

	err = client.Migration()
	assert.Nil(t, err)

	// create event error: foreign key constraint
	err = client.CreateEvent(Event{
		Password: "password",
		Time:     "time",
		Url:      "url",
	})
	assert.NotNil(t, err)

	// create user
	pwd, err := client.CreateUser("test@email.com")
	assert.Nil(t, err)

	// create event
	err = client.CreateEvent(Event{
		Password: pwd,
		Time:     time.DateTime,
		Url:      "https://demo.com",
	})
	assert.Nil(t, err)

	err = client.Close()
	assert.Nil(t, err)
}

func TestClient_FindEventByUser(t *testing.T) {
	defer tearDown()
	client, err := NewClient(rightDsn)
	assert.Nil(t, err)

	err = client.Migration()
	assert.Nil(t, err)

	// create users
	pwd1, _ := client.CreateUser("email1")
	pwd2, _ := client.CreateUser("email2")

	// create events
	for i := 0; i < 10; i++ {
		err = client.CreateEvent(Event{
			Password: pwd1,
			Time:     time.DateTime,
			Url:      fmt.Sprintf("https://demo.com/%d", i),
		})
		assert.Nil(t, err)
		err = client.CreateEvent(Event{
			Password: pwd2,
			Time:     time.DateTime,
			Url:      fmt.Sprintf("https://demo2.com/%d", i),
		})
		assert.Nil(t, err)
	}

	// find events
	events, err := client.FindEventByUser(pwd1)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(events))

	events, err = client.FindEventByUser(pwd2)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(events))

	err = client.Close()
}
