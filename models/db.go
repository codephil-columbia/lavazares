package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis"

	_ "github.com/lib/pq"
)

var db *sql.DB

//RedisCache is
var RedisCache *redis.Client

//InitDB initialized connected to database
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
		return err
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
		return err
	}

	return nil
}

//InitRedisCache initialzes the connection to the redis cache
func InitRedisCache() error {
	RedisCache = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := RedisCache.Ping().Result()
	fmt.Println(pong)
	if err != nil {
		log.Fatalln("error connecting to redis: %s", err)
		return err
	}
	return nil
}
