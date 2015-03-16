package collectors

import "strings"
import "syscall"
import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/Sirupsen/logrus"

type DiskspaceCollector struct{}

var filesystems = map[string]bool{
	"ext2":      true,
	"ext3":      true,
	"ext4":      true,
	"xfs":       true,
	"glusterfs": true,
	"rootfs":    true,
	"nfs":       true,
	"ntfs":      true,
	"hfs":       true,
	"fat32":     true,
	"fat16":     true,
	"btrfs":     true,
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
		if !filesystems[mount.FSType] {
			continue
		}

		syscall.Statfs(mount.MountPoint, &statfs_t)
		byte_avail := statfs_t.Bavail * uint64(statfs_t.Bsize)
		byte_free := statfs_t.Bfree * uint64(statfs_t.Bsize)

		diskspaceMapping[mount.Device] = mappings.MetricMap{
			"byte_avail": byte_avail,
			"byte_free":  byte_free,
			"byte_used":  byte_avail - byte_free,
			// "gigabyte_avail": statfs_t.Gigabyte_avail, // TODO
			// "gigabyte_free": statfs_t.Gigabyte_free, // TODO
			// "gigabyte_used": statfs_t.Gigabyte_used, // TODO
			// "inodes_avail": statfs_t.Inodes_avail // TODO
			// "inodes_free": statfs_t.Inodes_free // TODO
			// "inodes_used": statfs_t.Inodes_used // TODO
		}
	}

	return diskspaceMapping, nil
}

func (c *DiskspaceCollector) Report() (mappings.MetricMapSlice, error) {
	var report mappings.MetricMapSlice
	data, _ := c.Collect()

	if data != nil {
		units := map[string]string{
			"byte":   "B",
			"inodes": "Ino",
		}

		for device, values := range data {
			mountpoint := strings.Replace(device, "/", "_", -1)
			mountpoint = strings.Replace(mountpoint, ".", "_", -1)
			if mountpoint == "_" {
				mountpoint = "root"
			}

			for k, v := range values {
				s := strings.Split(k, "_")
				unit, mtype := s[0], s[1]
				report = append(report, mappings.MetricMap{
					"mountpoint":  mountpoint,
					"target_type": "gauge",
					"type":        mtype,
					"unit":        units[unit],
					"result":      v,
				})
			}
		}
	}

	return report, nil
}
