package connection

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

func ConnectDB() *sql.DB {
	dbClient, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_CONNECTION"))
	if nil != err {
		log.Fatal("open error:", err)
	}
	if err = dbClient.Ping(); err != nil {
		time.Sleep(time.Second * 3)
		fmt.Println("retry: connect to auth_db")
		return ConnectDB()
	}
	fmt.Println("connected to auth_db")
	return dbClient
}

var ctx = context.Background()

func ConnectRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_CONNECTION"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		time.Sleep(time.Second * 3)
		fmt.Println("retry: connect to redis")
		return ConnectRedis()
	}
	fmt.Println("connected to redis")
	return redisClient
}
