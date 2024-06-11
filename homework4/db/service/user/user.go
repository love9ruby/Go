package user

import (
	"db/model"
	"db/repository/mongo"
	"db/repository/orm"
	"db/repository/pg"
	rdb "db/repository/redis"
)

type UserServiceType int

const (
	Postgress UserServiceType = iota
	PgOrm                     // ORM
	Mongo
	Rdb // Redis
)

type UserService interface {
	InitTable() error
	GetType() UserServiceType
	CreateUser(email, name string) error
	UpdateUser(id, email string) error
	DeleteUser(id string) error
	GetUserList() ([]*model.User, error)
}

func NewUserService(t UserServiceType) *UserService {
	var ret UserService
	switch t {
	case Postgress:
		ret = new(PostgressUserService)
		ret.(*PostgressUserService).alias = t
		ret.(*PostgressUserService).db = pg.GetDB()
	case PgOrm:
		ret = new(PgOrmUserService)
		ret.(*PgOrmUserService).alias = t
		ret.(*PgOrmUserService).gorm = orm.GetDB()
	case Mongo:
		ret = new(MongoUserService)
		ret.(*MongoUserService).alias = t
		ret.(*MongoUserService).Client, ret.(*MongoUserService).DB = mongo.GetDB()
	case Rdb:
		ret = new(RedisUserService)
		ret.(*RedisUserService).alias = t
		ret.(*RedisUserService).client = rdb.GetDB()
	}
	return &ret
}

func Close(service UserService) {
	switch service.GetType() {
	case Postgress:
		pg.Close()
	case PgOrm:
		orm.Close()
	case Mongo:
		mongo.Close()
	case Rdb:
		rdb.Close()
	}
}
