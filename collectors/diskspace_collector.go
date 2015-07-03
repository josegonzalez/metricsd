package collectors

import "errors"
import "fmt"
import "os/exec"
import "strconv"
import "strings"
import "syscall"
import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type DiskspaceCollector struct {
	enabled        bool
	excludeFilters []string
	filesystems    map[string]bool
}

func (c *DiskspaceCollector) Enabled() bool {
	return c.enabled
}

func (c *DiskspaceCollector) State(state bool) {
	c.enabled = state
}

func (c *DiskspaceCollector) Setup(conf ini.File) {
	c.State(true)
	c.filesystems = map[string]bool{}

	fs, ok := conf.Get("DiskspaceCollector", "filesystems")
	if ok {
		enabledFilesystems := strings.Split(fs, ",")
		for _, enabledFilesystem := range enabledFilesystems {
			c.filesystems[strings.TrimSpace(enabledFilesystem)] = true
		}
	} else {
		c.filesystems["ext2"] = true
		c.filesystems["ext3"] = true
		c.filesystems["ext4"] = true
		c.filesystems["xfs"] = true
		c.filesystems["glusterfs"] = true
		c.filesystems["rootfs"] = true
		c.filesystems["nfs"] = true
		c.filesystems["ntfs"] = true
		c.filesystems["hfs"] = true
		c.filesystems["fat32"] = true
		c.filesystems["fat16"] = true
		c.filesystems["btrfs"] = true
	}

	ef, ok := conf.Get("DiskspaceCollector", "exclude_filters")
	if ok {
		excludeFilters := strings.Split(ef, ",")
		for _, excludeFilter := range excludeFilters {
			c.excludeFilters = append(c.excludeFilters, strings.TrimSpace(excludeFilter))
		}
	}
}

func (c *DiskspaceCollector) Collect() (map[string]mappings.MetricMap, error) {
	stat, err := linux.ReadMounts("/proc/mounts")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	var statfs_t syscall.Statfs_t
	diskspaceMapping := map[string]mappings.MetricMap{}

	for _, mount := range stat.Mounts {
		if !c.filesystems[mount.FSType] {
			continue
		}
		syscall.Statfs(mount.MountPoint, &statfs_t)
		byte_avail := statfs_t.Bavail * uint64(statfs_t.Bsize)
		byte_free := statfs_t.Bfree * uint64(statfs_t.Bsize)

		diskspaceMapping[mount.MountPoint] = mappings.MetricMap{
			"byte_avail":     byte_avail,
			"byte_free":      byte_free,
			"byte_used":      byte_avail - byte_free,
			"gigabyte_avail": byte_avail / 1073741824,
			"gigabyte_free":  byte_free / 1073741824,
			"gigabyte_used":  (byte_avail - byte_free) / 1073741824,
		}
	}

	dfMapping, e := collectDf()
	if e != nil {
		return diskspaceMapping, e
	}

	for mountpoint, metricMap := range dfMapping {
		if _, ok := diskspaceMapping[mountpoint]; !ok {
			continue
		}

		for key, value := range metricMap {
			diskspaceMapping[mountpoint][key] = value
		}
	}

	return diskspaceMapping, nil
}

func (c *DiskspaceCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	data, _ := c.Collect()

	if data != nil {
		units := map[string]string{
			"gigabyte": "GB",
			"byte":     "B",
			"inodes":   "Ino",
		}

		for mountpoint, values := range data {
			// TODO: Add exclude_filters support
			mountpoint = parseMountpoint(mountpoint)
			for k, v := range values {
				s := strings.Split(k, "_")
				unit, mtype := s[0], s[1]

				metric := structs.BuildMetric("diskspace", "gauge", mtype, v, structs.FieldsMap{
					"mountpoint": mountpoint,
					"unit":       units[unit],
					"raw_key":    k,
					"raw_value":  v,
				})
				metric.Path = fmt.Sprintf("diskspace.%s", mountpoint)
				report = append(report, metric)
			}
		}
	}

	return report, nil
}

func parseMountpoint(device string) string {
	mountpoint := strings.Replace(device, "/", "_", -1)
	mountpoint = strings.Replace(mountpoint, ".", "_", -1)
	if mountpoint == "_" {
		mountpoint = "root"
	}

	return mountpoint
}

func collectDf() (map[string]mappings.MetricMap, error) {
	data := map[string]mappings.MetricMap{}
	lines, e := readDf("-i")
	if e != nil {
		return data, e
	}

	for _, line := range lines[1:] {
		if !strings.HasPrefix(line, "/") {
			continue
		}
		chunks := strings.Fields(line)
		if len(chunks) >= 6 {
			mountpoint := chunks[5]
			if _, ok := data[mountpoint]; !ok {
				data[mountpoint] = mappings.MetricMap{}
			}

			if v, e := strconv.ParseInt(chunks[1], 10, 64); e == nil {
				data[mountpoint]["inodes_total"] = v
			}
			if v, e := strconv.ParseInt(chunks[2], 10, 64); e == nil {
				data[mountpoint]["inodes_used"] = v
			}
			if v, e := strconv.ParseInt(chunks[3], 10, 64); e == nil {
				data[mountpoint]["inodes_avail"] = v
			}
			if v, e := strconv.ParseInt(strings.Replace(chunks[4], "%", "", 1), 10, 64); e == nil {
				data[mountpoint]["inodes_use"] = v
			}
		}
	}
	return data, nil
}

func readDf(flag string) ([]string, error) {
	lines := []string{}
	raw, e := exec.Command("df", flag).Output()
	if e != nil {
		return lines, e
	}
	if len(raw) == 0 {
		return lines, errors.New("Reading df returned an empty string")
	}

	lines = strings.Split(strings.TrimSpace(string(raw)), "\n")
	return lines, nil
}
