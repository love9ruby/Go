# enable-redis-cluster

## How to create local redis cluster

> `en0` is your network interface that you're using right now.

```bash
ip=$(ipconfig getifaddr en0) docker-compose up -d --build
```

## How to connect to redis cluster

connect all ip

```go
// init cache
rp := new(UrlPoolCacheImpl)
rp.cache = rdb.NewClient([]string{"localhost:6999", "localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005"}, mp.db)
```

add password and use cluster syntax

```go
type Client struct {
	rdb *goredislib.ClusterClient
	rs  *redsync.Redsync
	mdb *mongo.Client
	ctx context.Context
}

func NewClient(addr []string, mdb *mongo.Client) *Client {
	client := goredislib.NewClusterClient(&goredislib.ClusterOptions{
		Addrs:    addr,
		Password: "pass.123", // no password set
	})
```