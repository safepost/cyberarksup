package main

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
	"os"
)

// disk usage of path/disk
func DiskUsage(letter string) bool {
	var free, total, avail uint64

	if len(letter) != 1 {
		logrus.Fatal("disk in configuration file must be a valid Windows Drive (eg : C or D)")
	}
	path := letter + ":\\"

	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		logrus.Fatal("disk in configuration file must be a valid Windows Drive")
	}

	//path := letter + ":\\"
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}
	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &avail)

	// fmt.Println(r1, r2, lastErr)
	logrus.Debug("Free:", free, "  Total:", total, "  Available:", avail)
	logrus.Debug((float64(avail) / float64(total)) * 100)

	// return True if free space is greater than 10%
	return (float64(avail)/float64(total))*100 > 10
}

// 	DiskUsage("c:\\")

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)
