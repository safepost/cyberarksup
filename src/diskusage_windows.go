//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

// disk usage of path/disk
func DiskUsage(letter string) bool {
	var free, total, avail uint64

	if len(letter) != 1 {
		log.Fatal("disk in configuration file must be a valid Windows Drive (eg : C or D)")
	}
	path := letter + ":\\"

	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		log.Panic("Disk in configuration file must be a valid Windows Drive " + path)
	}

	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}
	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &avail)
	if err != nil {
		log.Error("Unable to compute freespace for disk : " + letter)
		return false
	}

	percentageAvailable := (float64(avail) / float64(total)) * 100

	log.Debug(fmt.Sprintf("Drive %s: - Free: %s, Total: %s, Available: %s (%.1f%%)",
		letter,
		formatBytes(free),
		formatBytes(total),
		formatBytes(avail),
		percentageAvailable))

	// return True if free space is greater than 10%
	return percentageAvailable > 10
}
