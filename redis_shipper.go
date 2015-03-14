package main

import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"
import "github.com/josegonzalez/go-radixurl"

type RedisShipper struct{}

var redisList string
var redisUrl string

func (shipper *RedisShipper) Config(conf ini.File) {
	redisList = "metricsd"
	useRedisList, ok := conf.Get("RedisShipper", "list")
	if ok {
		redisList = useRedisList
	}

	redisUrl = "redis://127.0.0.1:6379/0"
	useRedisUrl, ok := conf.Get("RedisShipper", "url")
	if ok {
		redisUrl = useRedisUrl
	}
}

func (shipper *RedisShipper) Ship(logs MetricMapSlice) error {
	c, err := radixurl.ConnectToURL(redisUrl)
	errHndlr(err)
	defer c.Close()

	var list []string

	for _, item := range logs {
		serialized := MarshalData(item)
		list = append(list, string(serialized))
	}

	length := len(logs)
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
		logrus.Fatal("redis error: ", err)
	}
}
