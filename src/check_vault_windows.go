//go:build windows

package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func findVaultINIFile() ([]string, error) {
	initialSearchList := [...]string{"d:\\Cyberark", "d:\\Program Files\\Cyberark", "d:\\Program Files (x86)\\CyberArk",
		"c:\\Cyberark", "c:\\Program Files\\Cyberark", "c:\\Program Files (x86)\\CyberArk", "d:"}

	var finalSearchList []string
	for _, element := range initialSearchList {
		if stat, err := os.Stat(element); err == nil && stat.IsDir() {
			log.Debug(os.Stat(element))
			log.Debug(err)
			finalSearchList = append(finalSearchList, element)
			log.Debug("Keeping " + element)
		}
	}

	log.Debug("List of kept paths : " + strings.Join(finalSearchList, ","))

	if len(finalSearchList) == 0 {
		log.Info("Unable to find any suitable Vault.ini file in all known path! Specify it in config " +
			"file instead")
		os.Exit(1)
	}

	var fileList []string
	for _, validPath := range finalSearchList {
		e := filepath.Walk(validPath, func(path string, f os.FileInfo, err error) error {
			if strings.Contains(path, "Vault\\Vault.ini") {
				fileList = append(fileList, path)
			}
			return nil
		})

		if e != nil {
			panic(e)
		}
	}

	return fileList, nil
}
