//go:build linux

package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// disk usage of path/disk
func DiskUsage(path string) bool {

	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		log.Panic("Disk in configuration file must be a valid mount point " + path)
	}

	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		log.Error("Unable to get filesystem info for: " + path)
		return false
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)   // Espace libre total (incluant réservé)
	avail := stat.Bavail * uint64(stat.Bsize) // Espace disponible pour utilisateurs non-root
	used := total - free

	percentageUsed := (float64(used) / float64(total)) * 100
	percentageAvailable := (float64(avail) / float64(total)) * 100

	log.Debug(fmt.Sprintf("Path %s:", path))
	log.Debug(fmt.Sprintf("  Total: %s, Used: %s (%.1f%%), Available: %s (%.1f%%)",
		formatBytes(total),
		formatBytes(used),
		percentageUsed,
		formatBytes(avail),
		percentageAvailable))
	log.Debug(fmt.Sprintf("  Filesystem info - Blocks: %d, Block size: %d bytes, Free blocks: %d",
		stat.Blocks, stat.Bsize, stat.Bfree))

	// return True if free space is greater than 10%
	return percentageAvailable > 10
}
