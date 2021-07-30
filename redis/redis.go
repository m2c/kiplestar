package redis

import (
	"context"
	"errors"
	"fmt"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/m2c/kiplestar/config"
	"time"
)

type Redis struct {
	redisSource *redisv8.Client
	name        string //redis  name
}

func (slf *Redis) Redis() *redisv8.Client {
	return slf.redisSource
}

func (slf *Redis) Name() string {
	return slf.name
}

func (slf *Redis) StartRedis(config config.RedisConfig) error {
	if slf.redisSource != nil {
		return errors.New("redis already opened")
	}
	slf.name = config.Name
	readTimeout := 10 * time.Second
	slf.redisSource = redisv8.NewClient(&redisv8.Options{
		Addr:        config.Host,
		Password:    config.Password, // no password set
		DB:          config.Db,       // use default Client
		ReadTimeout: readTimeout,
	})
	timeout, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := slf.redisSource.Ping(timeout).Err()
	if err != nil {
		panic(fmt.Sprintf("redis connetc error %s", err.Error()))
	}
	return nil
}

func (slf *Redis) StopRedis() error {
	if slf.redisSource == nil {
		return errors.New("redis not opened")
	}
	err := slf.redisSource.Close()
	if err != nil {
		slf.redisSource = nil
	}
	return err
}
