package main

import "encoding/json"
import "fmt"
import "time"
import "github.com/fzzy/radix/redis"
import "github.com/Sirupsen/logrus"

type LogstashHook struct{}

func (hook *LogstashHook) Fire(entry *logrus.Entry) error {
    data := MarshalData(entry)

    serialized, err := json.Marshal(data)
    if err != nil {
        logrus.Error("Failed to marshal fields to JSON, %v", err)
        return nil
    }

    redisHost := Getenv("REDIS_HOST", "127.0.0.1")
    redisPort := Getenv("REDIS_PORT", "6379")

    c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort), time.Duration(10)*time.Second)
    errHndlr(err)
    defer c.Close()

    r := c.Cmd("rpush", "metricsd", serialized)
    errHndlr(r.Err)

    return nil
}

func (hook *LogstashHook) Levels() []logrus.Level {
  return []logrus.Level{
    logrus.InfoLevel,
  }
}
