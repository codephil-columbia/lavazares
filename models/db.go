package models

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

//RedisCache is
var redisCache *redis.Client

//InitDB initialized connected to database
func InitDB(dataSourceName string) error {
	var err error
	db, err = sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	return nil
}

func SetToSession(key string, val interface{}) error {
	return redisCache.Set(key, val, 0).Err()
}

func GetFromSession(key string) string {
	return redisCache.Get(key).Val()
}

func IsInSession(key string) int64 {
	return redisCache.Exists(key).Val()
}

func DeleteFromSession(key string) error {
	return redisCache.Del(key).Err()
}

func initTestDB(datasource string) (*gorm.DB, error) {
	testDB, err := gorm.Open("postgres", datasource)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	testDB.AutoMigrate(&User{})

	return testDB, err
}

//InitRedisCache initialzes the connection to the redis cache
func InitRedisCache() error {
	redisCache = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisCache.Ping().Result()
	fmt.Println(pong)
	if err != nil {
		log.Fatalln("error connecting to redis: %s", err)
		return err
	}
	return nil
}
