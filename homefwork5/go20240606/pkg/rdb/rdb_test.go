package rdb

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"short-url/pkg/mongo"
	"testing"
)

func setup() {
	mdb, err := mongo.NewClient("mongodb://localhost:27017", "test")
	if err != nil {
		panic(err)
		return
	}
	client = NewClient("localhost:6379", mdb)
}

func tearDown() {
	clearRedis()
	err := client.Close()
	if err != nil {
		return
	}
}

func clearRedis() {
	_, err := client.rdb.FlushAll(client.ctx).Result()
	if err != nil {
		return
	}
}

var client *Client

func TestNewClient(t *testing.T) {
	setup()
	defer tearDown()
}

func TestClient_SetUrl(t *testing.T) {
	setup()
	defer tearDown()
	url := Url{
		Password: "password",
		LongURL:  "long_url",
	}
	err := client.SetUrl("short_url", url)
	assert.Nil(t, err)

	url.Password = "password2"
	err = client.SetUrl("short_url2", url)
	assert.Nil(t, err)

	// test get
	url, err = client.GetUrl("short_url")
	assert.Nil(t, err)
	assert.Equal(t, "password", url.Password)
	assert.Equal(t, "long_url", url.LongURL)

	url, err = client.GetUrl("short_url2")
	assert.Nil(t, err)
	assert.Equal(t, "password2", url.Password)
	assert.Equal(t, "long_url", url.LongURL)
}

func TestClient_GetUrl(t *testing.T) {
	setup()
	defer tearDown()

	shortUrl := uuid.NewString()

	// not in cache, not in mongo
	url, err := client.GetUrl(shortUrl)
	assert.NotNil(t, err)
	assert.Equal(t, "", url.Password)
	assert.Equal(t, "", url.LongURL)

	// not in cache, but in mongo
	// push Url to mongo
	err = client.mdb.CreateUrl(mongo.Url{
		ShortURL: shortUrl,
		LongURL:  "long_url",
		Password: "password",
	})
	assert.Nil(t, err)
	url, err = client.GetUrl(shortUrl)
	assert.Nil(t, err)
	assert.Equal(t, "password", url.Password)
	assert.Equal(t, "long_url", url.LongURL)

	// in cache
	url, err = client.GetUrl(shortUrl)
	assert.Nil(t, err)
	assert.Equal(t, "password", url.Password)
	assert.Equal(t, "long_url", url.LongURL)
}
