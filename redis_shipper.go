package main

import "encoding/json"
import "fmt"
import "os"
import "time"
import "github.com/fzzy/radix/redis"
import log "github.com/Sirupsen/logrus"

type RedisShipper struct{}

func (hook *RedisShipper) Ship(logs []log.Fields) error {
	log.Debug(fmt.Sprintf("Shipping %d logs", len(logs)))

	redisHost := Getenv("REDIS_HOST", "127.0.0.1")
	redisPort := Getenv("REDIS_PORT", "6379")
	redisList := Getenv("REDIS_LIST", "metricsd")

	c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort), time.Duration(10)*time.Second)
	errHndlr(err)
	defer c.Close()

	var list []byte

	for _, item := range logs {
		serialized, err := json.Marshal(item)
		if err != nil {
			fmt.Errorf("Failed to marshal fields to JSON, %v", err)
			return nil
		}
		list = Extend(list, serialized)
	}

	r := c.Cmd("rpush", redisList, list)
	errHndlr(r.Err)

	return nil
}

func errHndlr(err error) {
	if err != nil {
		log.Error("error:", err)
		os.Exit(1)
	}
}
