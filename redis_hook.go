package main

import "encoding/json"
import "fmt"
import "time"
import "github.com/fzzy/radix/redis"
import log "github.com/Sirupsen/logrus"

type RedisHook struct{}

func (hook *RedisHook) Fire(entry *log.Entry) error {
	data := MarshalData(entry)

	serialized, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil
	}

	redisHost := Getenv("REDIS_HOST", "127.0.0.1")
	redisPort := Getenv("REDIS_PORT", "6379")
	redisList := Getenv("REDIS_LIST", "metricsd")

	c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort), time.Duration(10)*time.Second)
	errHndlr(err)
	defer c.Close()

	r := c.Cmd("rpush", redisList, serialized)
	errHndlr(r.Err)

	return nil
}

func (hook *RedisHook) Levels() []log.Level {
	return []log.Level{
		log.InfoLevel,
	}
}

func errHndlr(err error) {
	if err != nil {
		log.Error("error:", err)
		os.Exit(1)
	}
}
