package main

import "fmt"
import "os"
import "time"
import "github.com/fzzy/radix/redis"
import log "github.com/Sirupsen/logrus"

type RedisShipper struct{}

func (hook *RedisShipper) Ship(logs MetricMapSlice) error {
	length := len(logs)

	redisHost := Getenv("REDIS_HOST", "127.0.0.1")
	redisPort := Getenv("REDIS_PORT", "6379")
	redisList := Getenv("REDIS_LIST", "metricsd")

	c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort), time.Duration(10)*time.Second)
	errHndlr(err)
	defer c.Close()

	var list []string

	for _, item := range logs {
		serialized := MarshalData(item)
		list = append(list, string(serialized))
	}

	if length == 10 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2], list[3], list[4], list[5], list[6], list[7], list[8], list[9])
		errHndlr(r.Err)
	} else if length == 9 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2], list[3], list[4], list[5], list[6], list[7], list[8])
		errHndlr(r.Err)
	} else if length == 8 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2], list[3], list[4], list[5], list[6], list[7])
		errHndlr(r.Err)
	} else if length == 7 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2], list[3], list[4], list[5], list[6])
		errHndlr(r.Err)
	} else if length == 6 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2], list[3], list[4], list[5])
		errHndlr(r.Err)
	} else if length == 5 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2], list[3], list[4])
		errHndlr(r.Err)
	} else if length == 4 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2], list[3])
		errHndlr(r.Err)
	} else if length == 3 {
		r := c.Cmd("rpush", redisList, list[0], list[1], list[2])
		errHndlr(r.Err)
	} else if length == 2 {
		r := c.Cmd("rpush", redisList, list[0], list[1])
		errHndlr(r.Err)
	} else if length == 1 {
		r := c.Cmd("rpush", redisList, list[0])
		errHndlr(r.Err)
	}

	return nil
}

func errHndlr(err error) {
	if err != nil {
		log.Error("error:", err)
		os.Exit(1)
	}
}
