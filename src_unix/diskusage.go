package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
)

// disk usage of path/disk
func DiskUsage(path string) bool {
	var total, avail uint64

	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		log.Panic("Disk in configuration file must be a valid mount point " + path)
	}

	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		log.Error("Unable to get info from disk")
	}

	avail = stat.Bavail * uint64(stat.Bsize)
	total = stat.Bfree * uint64(stat.Bsize)
	// fmt.Println(r1, r2, lastErr)
	log.Debug((float64(avail) / float64(total)) * 100)

	// return True if free space is greater than 10%
	return (float64(avail)/float64(total))*100 > 10
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)
