package mq

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewIMQ(t *testing.T) {
	rMQ, err := NewIMQ()
	assert.Nil(t, err)
	assert.NotNil(t, rMQ)

	err = rMQ.Close()
	assert.Nil(t, err)
}

func TestNewRabbitMQ(t *testing.T) {
	rmq, err := NewRabbitMQ("userEvent", "", "")
	assert.Nil(t, err)
	assert.NotNil(t, rmq)

	err = rmq.Close()
	assert.Nil(t, err)
}

func TestRabbitMQ_Simple(t *testing.T) {
	rmq, err := NewRabbitMQ("userEvent", "", "")
	assert.Nil(t, err)
	assert.NotNil(t, rmq)

	msgChan, err := rmq.Consume()
	assert.Nil(t, err)
	assert.NotNil(t, msgChan)

	// try to fetch a message from the channel
	go func() {
		for i := 0; i < 10; i++ {
			j := map[string]string{"name": "test", "id": strconv.Itoa(i)}
			str, err := json.Marshal(j)
			assert.Nil(t, err)
			err = rmq.Publish(str)
			assert.Nil(t, err)
		}
	}()

	for i := 0; i < 10; i++ {
		data := <-msgChan
		assert.NotNil(t, data)
		assert.Equal(t, data.ContentType, "json/application")
		assert.NotNil(t, data.Body)
		// byte to map
		var j map[string]string
		err = json.Unmarshal(data.Body, &j)
		assert.Nil(t, err)
		assert.Equal(t, j["name"], "test")
		assert.Equal(t, j["id"], strconv.Itoa(i))
	}

	err = rmq.Close()
	assert.Nil(t, err)
}
