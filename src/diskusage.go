package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

// disk usage of path/disk
func DiskUsage(letter string) bool {
	var free, total, avail uint64

	if len(letter) != 1 {
		logrus.Fatal("disk in configuration file must be a valid Windows Drive (eg : C or D)")
	}

	path := letter + ":\\"
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}
	err = windows.GetDiskFreeSpaceEx(pathPtr, &free, &total, &avail)

	// fmt.Println(r1, r2, lastErr)
	fmt.Println("Free:", free, "Total:", total, "Available:", avail)
	fmt.Printf("%.2f", (float64(avail)/float64(total))*100)

	// return True if free space is greater than 10%
	return avail/total*100 > 10
}

// 	DiskUsage("c:\\")

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)
