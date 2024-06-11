package mongo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"reflect"
	"testing"
)

var client *Client

func dropDB(client *Client) error {
	err := client.DB.Drop(client.ctx)
	if err != nil {
		return err
	}
	return nil
}

func initClient() (*Client, error) {
	return NewClient("mongodb://localhost:27017", "test")
}

func setup() {
	println("Setting up")
	var err error
	client, err = initClient()
	if err != nil {
		fmt.Println("Error initializing client")
		os.Exit(1)
	}
}

func teardown() {
	println("Tearing down")
	err := dropDB(client)
	if err != nil {
		fmt.Println("Error dropping database")
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		url       string
		defaultDB []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "Test NewClient",
			args: args{
				url:       "mongodb://localhost:27017",
				defaultDB: []string{"test"},
			},
			want: &Client{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Test NewClient",
			args: args{
				url:       "mongodb://localhost:27016",
				defaultDB: []string{"test"},
			},
			want:    &Client{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.url, tt.args.defaultDB...)
			if tt.wantErr && err == nil {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.MongoClient == nil || got.DB == nil {
				t.Errorf("NewClient() = should not be nil")
			}
		})
	}
}

func TestClient_Ping(t *testing.T) {
	setup()
	defer teardown()
	type fields struct {
		MongoClient *mongo.Client
		DB          *mongo.Database
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test Ping",
			fields: fields{
				MongoClient: &mongo.Client{},
				DB:          &mongo.Database{},
				ctx:         context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Test Ping",
			fields: fields{
				MongoClient: nil,
				DB:          nil,
				ctx:         context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Test Ping",
			fields: fields{
				MongoClient: client.MongoClient,
				DB:          client.DB,
				ctx:         client.ctx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				MongoClient: tt.fields.MongoClient,
				DB:          tt.fields.DB,
				ctx:         tt.fields.ctx,
			}
			if err := c.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	setup()
	defer teardown()
	type fields struct {
		MongoClient *mongo.Client
		DB          *mongo.Database
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test Close",
			fields: fields{
				MongoClient: &mongo.Client{},
				DB:          &mongo.Database{},
				ctx:         context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Test Close",
			fields: fields{
				MongoClient: nil,
				DB:          nil,
				ctx:         context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Test Close",
			fields: fields{
				MongoClient: client.MongoClient,
				DB:          client.DB,
				ctx:         client.ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				MongoClient: tt.fields.MongoClient,
				DB:          tt.fields.DB,
				ctx:         tt.fields.ctx,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateUrl(t *testing.T) {
	setup()
	defer teardown()
	type fields struct {
		MongoClient *mongo.Client
		DB          *mongo.Database
		ctx         context.Context
	}
	type args struct {
		url Url
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test CreateUrl",
			fields: fields{
				MongoClient: client.MongoClient,
				DB:          client.DB,
				ctx:         client.ctx,
			},
			args: args{
				url: Url{
					ShortURL: "short",
					LongURL:  "long",
					Password: "password",
				},
			},
			wantErr: false,
		}, {
			name: "Test CreateUrl",
			fields: fields{
				MongoClient: client.MongoClient,
				DB:          client.DB,
				ctx:         client.ctx,
			},
			args: args{
				url: Url{
					ShortURL: "short",
					LongURL:  "long",
					Password: "password",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				MongoClient: tt.fields.MongoClient,
				DB:          tt.fields.DB,
				ctx:         tt.fields.ctx,
			}
			if err := c.CreateUrl(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("CreateUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_FindUrl(t *testing.T) {
	setup()
	//defer teardown()
	// Insert a URL
	_ = client.CreateUrl(Url{
		ShortURL: "short",
		LongURL:  "long",
		Password: "password",
	})
	// end of insertion
	type fields struct {
		MongoClient *mongo.Client
		DB          *mongo.Database
		ctx         context.Context
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Url
		wantErr bool
	}{
		{
			name: "Test FindUrl",
			fields: fields{
				MongoClient: client.MongoClient,
				DB:          client.DB,
				ctx:         client.ctx,
			},
			args: args{
				shortURL: "shortNotExists",
			},
			want:    Url{},
			wantErr: true,
		},
		{
			name: "Test FindUrl",
			fields: fields{
				MongoClient: client.MongoClient,
				DB:          client.DB,
				ctx:         client.ctx,
			},
			args: args{
				shortURL: "short",
			},
			want: Url{
				ShortURL: "short",
				LongURL:  "long",
				Password: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				MongoClient: tt.fields.MongoClient,
				DB:          tt.fields.DB,
				ctx:         tt.fields.ctx,
			}
			got, err := c.FindUrl(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID == "" && !tt.wantErr {
				t.Errorf("object ID is empty")
			}
			got.ID = ""
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_InitKeyIndex(t *testing.T) {
	setup()
	defer teardown()

	_ = client.CreateKey("key")

	// create index
	err := client.InitKeyIndex()
	assert.Nil(t, err)

	// set not unique key
	err = client.CreateKey("key")
	assert.NotNil(t, err)

	// set unique key
	err = client.CreateKey("key2")
	assert.Nil(t, err)

	// drop the db
	_ = dropDB(client)

	// create 2 same keys
	err = client.CreateKey("key")
	assert.Nil(t, err)
	err = client.CreateKey("key")
	assert.Nil(t, err)

	// create index error
	err = client.InitKeyIndex()
	assert.NotNil(t, err)
}

func TestClient_IsKeyExist(t *testing.T) {
	setup()
	defer teardown()
	_ = client.InitKeyIndex()
	_ = client.CreateKey("key")
	_ = client.CreateKey("key2")
	_ = client.CreateKey("key3")

	isExist, err := client.IsKeyExist("key")
	assert.Nil(t, err)
	assert.True(t, isExist)

	isExist, err = client.IsKeyExist("key4")
	assert.Nil(t, err)
	assert.False(t, isExist)
}
