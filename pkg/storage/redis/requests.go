package redis

import (
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"log"
)

type Redis struct {
	Redis *redis.Client
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func (r *Redis) Save(chatID, taskID string) error {
	opt, err := redis.ParseURL("redis://my_user@localhost:6380/0")
	if err != nil {
		return err
	}

	r.Redis = redis.NewClient(opt)

	err = r.Redis.Set(chatID, taskID, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Get(chatID string) (error, string) {
	opt, err := redis.ParseURL("redis://my_user@localhost:6380/0")
	if err != nil {
		return err, ""
	}

	r.Redis = redis.NewClient(opt)

	taskID, err := r.Redis.Get(chatID).Result()
	if err != nil {
		return err, ""
	}

	return nil, taskID
}
