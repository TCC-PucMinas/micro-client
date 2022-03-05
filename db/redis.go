package db

import (
	"github.com/go-redis/redis"
)

var instanceDBRedis *redis.Client

func ConnectDatabaseRedis() *redis.Client {

	if instanceDBRedis != nil {
		if _, err := instanceDBRedis.Ping().Result(); err != nil {
			instanceDBRedis = nil
		}
	}

	if instanceDBRedis == nil {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", //os.Getenv("DB_REDIS_HOST"),     //
			Password: "SUASENHA",       //os.Getenv("DB_REDIS_PASSWORD"), //
			DB:       0,
		})

		if _, err := client.Ping().Result(); err != nil {
			panic(err.Error())
		}

		instanceDBRedis = client
	}

	return instanceDBRedis
}

func RemoveCacheRedisByKey(key string) {
	db := ConnectDatabaseRedis()
	db.Del(key)
}
