package collectors

import "fmt"
import "regexp"
import "strings"
import "github.com/josegonzalez/go-radixurl"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/josegonzalez/metricsd/utils"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type RedisCollector struct {
	enabled bool
	url     string
}

func (this *RedisCollector) Enabled() bool {
	return this.enabled
}

func (this *RedisCollector) State(state bool) {
	this.enabled = state
}

func (this *RedisCollector) Setup(conf ini.File) {
	this.State(true)

	useRedisUrl, ok := conf.Get("RedisCollector", "url")
	if ok {
		this.url = useRedisUrl
	} else {
		this.url = "redis://127.0.0.1:6379/0"
	}
}

func (this *RedisCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	data, _ := this.collect()

	if data != nil {
		// TODO: _ is a prefix
		for _, values := range data {
			for k, v := range values {
				metric := structs.BuildMetric("RedisCollector", "redis", "gauge", k, v, structs.FieldsMap{
					"raw_key":   k,
					"raw_value": v,
				})
				report = append(report, metric)
			}
		}
	}

	return report, nil
}

func (this *RedisCollector) collect() (map[string]mappings.MetricMap, error) {
	c, err := radixurl.ConnectToURL(this.url)
	errHndlr(err)
	defer c.Close()

	r := c.Cmd("info")
	errHndlr(r.Err)

	s, err := r.Str()
	errHndlr(err)

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
			redisMapping[fmt.Sprintf("db%d", res[1])] = mappings.MetricMap{
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

	// redisMapping["UptimeInSeconds", parseInt64(values["uptime_in_seconds"]))
	redisMapping["clients"] = mappings.MetricMap{
		"biggest_input_buf":   values["client_biggest_input_buf"],
		"blocked":             values["blocked_clients"],
		"connected":           values["connected_clients"],
		"longest_output_list": values["client_longest_output_list"],
	}

	// TODO
	// 'cpu.parent.sys': 'used_cpu_sys',
	// 'cpu.children.sys': 'used_cpu_sys_children',
	// 'cpu.parent.user': 'used_cpu_user',
	// 'cpu.children.user': 'used_cpu_user_children',

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

	return redisMapping, nil
}

func errHndlr(err error) {
	if err != nil {
		logrus.Fatal("redis error: ", err)
	}
}
