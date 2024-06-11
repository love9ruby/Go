package pool

import (
	"errors"
	"fmt"
	"go20240606/pkg/key"
	"go20240606/pkg/mongo"
	"go20240606/pkg/rdb"
	"os"
	"sync"
	"time"
)

type IUrlPool interface {
	// GetLongURL get the long URL from the short URL
	GetLongURL(shortURL string) (mongo.Url, error)
	// CreateShortURL create a short URL from the long URL
	CreateShortURL(longURL, password string) (string, error)
}

type UrlPool struct {
	db    *UrlPoolDBImpl
	cache *UrlPoolCacheImpl
	mtx   *sync.Mutex
	keys  []string
}

type UrlPoolDBImpl struct {
	db *mongo.Client
}

type UrlPoolCacheImpl struct {
	cache *rdb.Client
}

func (up *UrlPoolCacheImpl) GetUrl(shortURL string) (rdb.Url, error) {
	// get the long URL from the cache, cache miss will hit DB
	url, err := up.cache.GetUrl(shortURL)
	if err != nil {
		return rdb.Url{}, err
	}
	return url, nil
}

func (up *UrlPoolDBImpl) CreateUrl(shortURL, longURL, password string) error {
	url := mongo.Url{
		ShortURL: shortURL,
		LongURL:  longURL,
		Password: password,
	}
	return up.db.CreateUrl(url)
}

func (up *UrlPoolDBImpl) CreateNewTokens(num int) ([]string, error) {
	seed := int64(time.Now().Nanosecond())
	keys := key.GenerateKeys(num, seed, up.db.CreateKey)
	if len(keys) < 1 {
		return nil, errors.New("too less new key")
	}
	return keys, nil
}

func (up *UrlPool) GetLongURL(shortURL string) (mongo.Url, error) {
	// check cache, cache miss will hit DB
	rdbUrl, err := up.cache.GetUrl(shortURL)
	if err != nil {
		return mongo.Url{}, err
	}
	return mongo.Url{
		ShortURL: shortURL,
		LongURL:  rdbUrl.LongURL,
		Password: rdbUrl.Password,
	}, nil
}

func (up *UrlPool) CreateShortURL(longURL, password string) (string, error) {
	// get new token
	var shortURL string
	up.mtx.Lock()
	if len(up.keys) == 0 {
		newKeys, err := up.db.CreateNewTokens(10)
		if err != nil {
			return "", err
		}
		up.keys = append(up.keys, newKeys...)
	}
	shortURL = up.keys[0]
	up.keys = up.keys[1:]
	up.mtx.Unlock()
	// save new Url
	err := up.db.CreateUrl(shortURL, longURL, password)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func NewUrlPool() (pool IUrlPool, err error) {
	// env
	dbHost := os.Getenv("MONGO_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	cacheHost := os.Getenv("REDIS_HOST")
	if cacheHost == "" {
		cacheHost = "localhost"
	}
	// init DB
	mp := new(UrlPoolDBImpl)
	mp.db, err = mongo.NewClient(fmt.Sprintf("mongodb://%s:27017", dbHost))
	if err != nil {
		return nil, err
	}
	// init cache
	rp := new(UrlPoolCacheImpl)
	rp.cache = rdb.NewClient(fmt.Sprintf("%s:6379", cacheHost), mp.db)
	// wrap the pool
	pool = &UrlPool{
		db:    mp,
		cache: rp,
		mtx:   &sync.Mutex{},
		keys:  make([]string, 0),
	}

	// migrate the DB
	err = mp.db.InitKeyIndex()
	if err != nil {
		return nil, err
	}
	return pool, nil
}
