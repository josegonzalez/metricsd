package collectors

import "fmt"
import "regexp"
import "strings"
import "github.com/josegonzalez/go-radixurl"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/josegonzalez/metricsd/utils"
import "github.com/vaughan0/go-ini"

type RedisCollector struct {
	enabled bool
	url     string
}

func (c *RedisCollector) Enabled() bool {
	return c.enabled
}

func (c *RedisCollector) State(state bool) {
	c.enabled = state
}

func (c *RedisCollector) Setup(conf ini.File) {
	c.State(true)

	useRedisURL, ok := conf.Get("RedisCollector", "url")
	if ok {
		c.url = useRedisURL
	} else {
		c.url = "redis://127.0.0.1:6379/0"
	}
}

func (c *RedisCollector) Report() (structs.MetricSlice, error) {
	return c.collect()
}

func (c *RedisCollector) collect() (structs.MetricSlice, error) {
	conn, err := radixurl.ConnectToURL(c.url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	r := conn.Cmd("info")
	if r.Err != nil {
		return nil, r.Err
	}

	s, err := r.Str()
	if err != nil {
		return nil, err
	}

	dbRegexp := regexp.MustCompile("^db(\\d+):keys=(\\d+),expires=(\\d+)")

	redisMapping := map[string]mappings.MetricMap{}
	values := map[string]string{}

	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || len(line) < 2 {
			continue
		}
		res := dbRegexp.FindStringSubmatch(line)
		if len(res) == 4 {
			redisMapping[fmt.Sprintf("db%s", res[1])] = mappings.MetricMap{
				"expires": utils.ParseInt64(res[2]),
				"keys":    utils.ParseInt64(res[3]),
			}
			continue
		}
		chunks := strings.Split(line, ":")
		if len(chunks) > 1 {
			values[chunks[0]] = strings.Join(chunks[1:], ":")
		}
	}

	redisMapping = c.processValues(redisMapping, values)

	var report structs.MetricSlice
	for prefix, values := range redisMapping {
		for k, v := range values {
			metric := structs.BuildMetric("RedisCollector", "redis", "gauge", k, v, structs.FieldsMap{
				"raw_key":   k,
				"raw_value": v,
			})
			metric.Path = fmt.Sprintf("redis.%s", prefix)
			report = append(report, metric)
		}
	}

	metric := structs.BuildMetric("RedisCollector", "redis", "gauge", "sys", values["used_cpu_sys"], structs.FieldsMap{
		"raw_key":   "sys",
		"raw_value": values["used_cpu_sys"],
	})
	metric.Path = "redis.cpu.parent"
	report = append(report, metric)

	metric = structs.BuildMetric("RedisCollector", "redis", "gauge", "sys", values["used_cpu_sys_children"], structs.FieldsMap{
		"raw_key":   "sys",
		"raw_value": values["used_cpu_sys_children"],
	})
	metric.Path = "redis.cpu.children"
	report = append(report, metric)

	metric = structs.BuildMetric("RedisCollector", "redis", "gauge", "user", values["used_cpu_user"], structs.FieldsMap{
		"raw_key":   "user",
		"raw_value": values["used_cpu_user"],
	})
	metric.Path = "redis.cpu.parent"
	report = append(report, metric)

	metric = structs.BuildMetric("RedisCollector", "redis", "gauge", "user", values["used_cpu_user_children"], structs.FieldsMap{
		"raw_key":   "user",
		"raw_value": values["used_cpu_user_children"],
	})
	metric.Path = "redis.cpu.children"
	report = append(report, metric)

	return report, nil
}

func (c *RedisCollector) processValues(redisMapping map[string]mappings.MetricMap, values map[string]string) map[string]mappings.MetricMap {
	redisMapping["clients"] = mappings.MetricMap{
		"biggest_input_buf":   values["client_biggest_input_buf"],
		"blocked":             values["blocked_clients"],
		"connected":           values["connected_clients"],
		"longest_output_list": values["client_longest_output_list"],
	}

	redisMapping["hash_max_zipmap"] = mappings.MetricMap{
		"entries": values["hash_max_zipmap_entries"],
		"value":   values["hash_max_zipmap_value"],
	}
	redisMapping["keys"] = mappings.MetricMap{
		"evicted": values["evicted_keys"],
		"expired": values["expired_keys"],
	}
	redisMapping["keyspace"] = mappings.MetricMap{
		"hits":   values["keyspace_hits"],
		"misses": values["keyspace_misses"],
	}
	redisMapping["last_save"] = mappings.MetricMap{
		"changes_since": values["changes_since_last_save"],
		"time":          values["last_save_time"],
	}
	redisMapping["memory"] = mappings.MetricMap{
		"internal_view":       values["used_memory"],
		"external_view":       values["used_memory_rss"],
		"fragmentation_ratio": values["mem_fragmentation_ratio"],
	}
	redisMapping["process"] = mappings.MetricMap{
		"commands_processed":   values["total_commands_processed"],
		"connections_received": values["total_connections_received"],
		"uptime":               values["uptime_in_seconds"],
	}
	redisMapping["pubsub"] = mappings.MetricMap{
		"channels": values["pubsub_channels"],
		"patterns": values["pubsub_patterns"],
	}
	redisMapping["slaves"] = mappings.MetricMap{
		"connected": values["connected_slaves"],
		"last_io":   values["master_last_io_seconds_ago"],
	}

	return redisMapping
}
