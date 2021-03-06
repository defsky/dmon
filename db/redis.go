package db

import (
	"fmt"
	"log"
	"time"

	"github.com/defsky/dmon/config"
	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func initRedis() {
	cfg := config.GetConfig().DB.Redis
	if cfg == nil {
		return
	}

	log.Println("Init Redis  ...")

	server := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	log.Printf("Connecting %s ...", server)

	redisPool = &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server,
				redis.DialConnectTimeout(30*time.Second),
				redis.DialReadTimeout(60*time.Second),
				redis.DialWriteTimeout(60*time.Second),
			)
			if err != nil {
				log.Printf("redis error: %s\n", err.Error())
				time.Sleep(10 * time.Second)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
