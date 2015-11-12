package collectors

import (
	"encoding/json"
	"reflect"
	"strconv"
)

import "fmt"
import "io/ioutil"
import "net/http"
import "strings"
import "github.com/Sirupsen/logrus"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

// ElasticsearchCollector is an exported type that
// allows collecting metrics for Elasticsearch
type ElasticsearchCollector struct {
	enabled   bool
	instances []string
}

// Enabled allows checking whether the collector is enabled or not
func (c *ElasticsearchCollector) Enabled() bool {
	return c.enabled
}

// State allows setting the enabled state of the collector
func (c *ElasticsearchCollector) State(state bool) {
	c.enabled = state
}

// Setup configures the collector
func (c *ElasticsearchCollector) Setup(conf ini.File) {
	c.State(true)
	instances, ok := conf.Get("ElasticsearchCollector", "instances")
	if !ok {
		instances = "http://127.0.0.1:9200"
	}
	for _, instance := range strings.Split(instances, ",") {
		c.instances = append(c.instances, instance)
	}
}

// Report collects a list of MetricSlices for upstream reporting
func (c *ElasticsearchCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	for _, instance := range c.instances {
		values, err := c.collect(instance, report)
		if err != nil {
			logrus.Warning(err)
		}
		for _, metric := range values {
			report = append(report, metric)
		}
	}
	return report, nil
}

func (c *ElasticsearchCollector) collect(instance string, report structs.MetricSlice) (structs.MetricSlice, error) {
	report = collectInstance(instance, report)
	report = collectInstanceClusterStats(instance, report)
	// TODO: Collect index stats
	return report, nil
}

func collectInstance(instance string, report structs.MetricSlice) structs.MetricSlice {
	url := fmt.Sprintf("%s/_nodes/_local/stats?all=true", instance)
	var result map[string]interface{}
	response, err := request(url)
	if err != nil {
		logrus.Warning(err)
		return report
	}

	err = json.Unmarshal(response, &result)
	if err != nil {
		logrus.Warning(err)
		return report
	}

	nodesData := result["nodes"].(map[string]interface{})
	for _, nodesData := range nodesData {
		data := nodesData.(map[string]interface{})
		report = appendMetric(report, data, "http.current", "http.current_open")

		if data["indices"] != nil {
			report = collectIndexStats(report, data)
		}

		if data["store"] != nil {
			report = appendMetric(report, data, "indices.datastore.size", "store.size_in_bytes")
		}

		if data["transport"] != nil {
			report = collectTransportStats(report, data)
		}

		if data["jvm"] != nil {
			report = collectJvmStats(report, data)
		}

		if data["thread_pool"] != nil {
			report = collectThreadPoolStats(report, data)
		}

		if data["network"] != nil {
			report = collectNetworkStats(report, data)
		}
	}

	return report
}

func collectIndexStats(report structs.MetricSlice, data map[string]interface{}) structs.MetricSlice {
	report = appendMetric(report, data, "indices.docs.count", "indices.docs.count")
	report = appendMetric(report, data, "indices.docs.deleted", "indices.docs.deleted")

	// elasticsearch < 0.90RC2
	report = appendMetric(report, data, "cache.bloom.size", "indices.cache.bloom_size_in_bytes")
	report = appendMetric(report, data, "cache.field.evictions", "indices.cache.field_evictions")
	report = appendMetric(report, data, "cache.field.size", "indices.cache.field_size_in_bytes")
	report = appendMetric(report, data, "cache.filter.count", "indices.cache.filter_count")
	report = appendMetric(report, data, "cache.filter.evictions", "indices.cache.filter_evictions")
	report = appendMetric(report, data, "cache.filter.size", "indices.cache.filter_size_in_bytes")
	report = appendMetric(report, data, "cache.id.size", "indices.cache.id_cache_size_in_bytes")

	// elasticsearch >= 0.90RC2
	report = appendMetric(report, data, "cache.filter.evictions", "indices.filter_cache.evictions")
	report = appendMetric(report, data, "cache.filter.size", "indices.filter_cache.memory_size_in_bytes")
	report = appendMetric(report, data, "cache.filter.count", "indices.filter_cache.count")

	// elasticsearch >= 0.90RC2
	report = appendMetric(report, data, "cache.id.size", "indices.id_cache.memory_size_in_bytes")

	// elasticsearch >= 0.90
	report = appendMetric(report, data, "fielddata.size", "indices.fielddata.memory_size_in_bytes")
	report = appendMetric(report, data, "fielddata.evictions", "indices.fielddata.evictions")

	// process mem/cpu (may not be present, depending on access restrictions)
	report = appendMetric(report, data, "process.cpu.percent", "process.cpu.percent")
	report = appendMetric(report, data, "process.mem.resident", "process.mem.resident_in_bytes")
	report = appendMetric(report, data, "process.mem.share", "process.mem.share_in_bytes")
	report = appendMetric(report, data, "process.mem.virtual", "process.mem.total_virtual_in_bytes")

	report = appendMetric(report, data, "disk.reads.count", "fs.data.0.disk_reads")
	report = appendMetric(report, data, "disk.reads.size", "fs.data.0.disk_read_size_in_bytes")
	report = appendMetric(report, data, "disk.writes.count", "fs.data.0.disk_writes")
	report = appendMetric(report, data, "disk.writes.size", "fs.data.0.disk_write_size_in_bytes")
	return report
}

func collectTransportStats(report structs.MetricSlice, data map[string]interface{}) structs.MetricSlice {
	report = appendMetric(report, data, "transport.rx.count", "transport.rx_count")
	report = appendMetric(report, data, "transport.rx.size", "transport.rx_size_in_bytes")
	report = appendMetric(report, data, "transport.tx.count", "transport.tx_count")
	report = appendMetric(report, data, "transport.tx.size", "transport.tx_size_in_bytes")
	return report
}

func collectJvmStats(report structs.MetricSlice, data map[string]interface{}) structs.MetricSlice {
	report = appendMetric(report, data, "jvm.mem.heap_used", "jvm.mem.heap_used_in_bytes")
	report = appendMetric(report, data, "jvm.mem.heap_committed", "jvm.mem.heap_committed_in_bytes")
	report = appendMetric(report, data, "jvm.mem.non_heap_used", "jvm.mem.non_heap_used_in_bytes")
	report = appendMetric(report, data, "jvm.mem.non_heap_committed", "jvm.mem.non_heap_committed_in_bytes")
	report = appendMetric(report, data, "jvm.mem.heap_used_percent", "jvm.mem.heap_used_percent")

	poolMapping := make(map[string]string)
	jvm := data["jvm"].(map[string]interface{})
	mem := jvm["mem"].(map[string]interface{})
	pools := mem["pools"].(map[string]interface{})
	gc := jvm["gc"].(map[string]interface{})
	collectors := gc["collectors"].(map[string]interface{})

	for pool := range pools {
		poolMapping[pool] = strings.Replace(pool, " ", "_", -1)
	}

	for lookup, replacement := range poolMapping {
		report = appendMetric(report, data, fmt.Sprintf("jvm.mem.pools.%s.used", replacement), fmt.Sprintf("jvm.mem.pools.%s.used_in_bytes", lookup))
		report = appendMetric(report, data, fmt.Sprintf("jvm.mem.pools.%s.max", replacement), fmt.Sprintf("jvm.mem.pools.%s.max_in_bytes", lookup))
	}

	report = appendMetric(report, data, "jvm.threads.count", "jvm.threads.count")

	collectionCount := 0.0
	collectionTimeInMillis := 0.0
	for collector, d := range collectors {
		report = appendMetric(report, data, fmt.Sprintf("jvm.gc.collection.%s.count", collector), fmt.Sprintf("jvm.gc.collectors.%s.collection_count", collector))
		report = appendMetric(report, data, fmt.Sprintf("jvm.gc.collection.%s.time", collector), fmt.Sprintf("jvm.gc.collectors.%s.collection_time_in_millis", collector))
		collectorData := d.(map[string]interface{})
		collectionCount += collectorData["collection_count"].(float64)
		collectionTimeInMillis += collectorData["collection_time_in_millis"].(float64)
	}

	// calculate the totals, as they're absent in elasticsearch >
	// 0.90.10
	if gc["collection_count"] != nil {
		report = appendMetric(report, data, "jvm.gc.collection.count", "jvm.gc.collection_count")
	} else {
		metric := structs.BuildMetric("ElasticsearchCollector", "elasticsearch", "gauge", "jvm.gc.collection.count", collectionCount, structs.FieldsMap{
			"raw_key":   nil,
			"raw_value": collectionCount,
		})
		report = append(report, metric)
	}

	if gc["collection_time_in_millis"] != nil {
		report = appendMetric(report, data, "jvm.gc.collection.time", "jvm.gc.collection_time_in_millis")
	} else {
		metric := structs.BuildMetric("ElasticsearchCollector", "elasticsearch", "gauge", "jvm.gc.collection.time", collectionTimeInMillis, structs.FieldsMap{
			"raw_key":   nil,
			"raw_value": collectionTimeInMillis,
		})
		report = append(report, metric)
	}

	return report
}

func collectThreadPoolStats(report structs.MetricSlice, data map[string]interface{}) structs.MetricSlice {
	for key, threadPoolData := range data["thread_pool"].(map[string]interface{}) {
		for k := range threadPoolData.(map[string]interface{}) {
			lookup := fmt.Sprintf("thread_pool.%s.%s", key, k)
			report = appendMetric(report, data, lookup, lookup)
		}
	}
	return report
}

func collectNetworkStats(report structs.MetricSlice, data map[string]interface{}) structs.MetricSlice {
	for key, threadPoolData := range data["network"].(map[string]interface{}) {
		for k := range threadPoolData.(map[string]interface{}) {
			lookup := fmt.Sprintf("network.%s.%s", key, k)
			report = appendMetric(report, data, lookup, lookup)
		}
	}
	return report
}

func collectInstanceClusterStats(instance string, report structs.MetricSlice) structs.MetricSlice {
	url := fmt.Sprintf("%s/_cluster/health", instance)
	var result map[string]interface{}
	response, err := request(url)
	if err != nil {
		logrus.Warning(err)
		return report
	}

	err = json.Unmarshal(response, &result)
	if err != nil {
		logrus.Warning(err)
		return report
	}

	var healthMap = map[string]string{
		"number_of_nodes":           "cluster_health.nodes.total",
		"number_of_data_nodes":      "cluster_health.nodes.data",
		"active_primary_shards":     "cluster_health.shards.active_primary",
		"active_shards":             "cluster_health.shards.active",
		"relocating_shards":         "cluster_health.shards.relocating",
		"unassigned_shards":         "cluster_health.shards.unassigned",
		"initializing_shards":       "cluster_health.shards.initializing",
		"delayed_unassigned_shards": "cluster_health.shards.delayed_unassigned",
		"number_of_pending_tasks":   "cluster_health.tasks.pending",
		"number_of_in_flight_fetch": "cluster_health.shards.info_requests",
	}
	for key, name := range healthMap {
		value, ok := result[key]
		if !ok {
			continue
		}
		metric := structs.BuildMetric("ElasticsearchCollector", "elasticsearch", "gauge", name, value, structs.FieldsMap{
			"raw_key":   key,
			"raw_value": value,
		})
		report = append(report, metric)
	}

	return report
}

func appendMetric(report structs.MetricSlice, data map[string]interface{}, name string, rawKey string) structs.MetricSlice {
	var lookupValue interface{}
	var value interface{}
	var ok bool
	lookupValue = data
	for _, lookup := range strings.Split(rawKey, ".") {
		original := reflect.ValueOf(lookupValue)
		if original.Kind() == reflect.Slice {
			lookupInt, err := strconv.Atoi(lookup)
			if err != nil {
				return report
			}
			lookupArray := lookupValue.([]interface{})
			if len(lookupArray) < lookupInt {
				return report
			}
			lookupValue = lookupArray[lookupInt]
			if !ok {
				return report
			}
			value = lookupValue
		} else if original.Kind() == reflect.Map {
			lookupValue, ok = lookupValue.(map[string]interface{})[lookup]
			if !ok {
				return report
			}
			value = lookupValue
		} else {
			return report
		}
	}

	metric := structs.BuildMetric("ElasticsearchCollector", "elasticsearch", "gauge", name, value, structs.FieldsMap{
		"raw_key":   rawKey,
		"raw_value": value,
	})
	report = append(report, metric)
	return report
}

func request(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
