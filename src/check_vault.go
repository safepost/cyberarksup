package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func testConnection(addr string, portOptional ...int) error {
	port := 1858
	if len(portOptional) > 0 {
		port = portOptional[0]
	}

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.DialContext(ctx, "tcp", addr+":"+strconv.Itoa(port))

	return err
}

// Find ADDRESS filed in Vault.ini file and return list of IP Address (string format)
func findVaultIPAddress(iniFilePath string) ([]string, error) {
	cfg, err := ini.Load(iniFilePath)
	if err != nil {
		logrus.Fatal("Unable to load ini file provided in configuration file : " + iniFilePath)
	}

	vaultIPAddresses := cfg.Section("").Key("ADDRESS").String()

	return strings.Split(vaultIPAddresses, ","), nil

}

func findVaultINIFile() ([]string, error) {
	initialSearchList := [...]string{"d:\\Cyberark", "d:\\Program Files\\Cyberark", "d:\\Program Files (x86)\\CyberArk",
		"c:\\Cyberark", "c:\\Program Files\\Cyberark", "c:\\Program Files (x86)\\CyberArk", "d:"}

	var finalSearchList []string
	for _, element := range initialSearchList {
		if stat, err := os.Stat(element); err == nil && stat.IsDir() {
			logrus.Debug(os.Stat(element))
			logrus.Debug(err)
			finalSearchList = append(finalSearchList, element)
			logrus.Debug("Keeping " + element)
		}
	}

	logrus.Debug("List of kept paths : " + strings.Join(finalSearchList, ","))

	if len(finalSearchList) == 0 {
		logrus.Fatal("Unable to find any suitable Vault.ini file in all known path! Specify it in config " +
			"file instead")
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

func getVaultsIPs(iniFile string) []string {
	var vaultIPs []string

	if iniFile == "" {
		logrus.Info("No file provided in configuration file, trying to find Vault.ini file")
		list, err := findVaultINIFile()

		if err != nil {
			logrus.Fatal("Unable to find Vault INI file")
		}

		logrus.Info("Using Vault.ini file : " + list[0])

		vaultIPs, _ = findVaultIPAddress(list[0])
	} else {
		vaultIPs, _ = findVaultIPAddress(iniFile)
	}

	logrus.Info("Using addresses : " + strings.Join(vaultIPs, ","))
	return vaultIPs
}

func checkVault(vaultIPs []string) bool {

	for _, ipAddr := range vaultIPs {
		logrus.Debug("Testing connection to ", ipAddr)
		err := testConnection(ipAddr)
		if err != nil {
			logrus.Debug("Vault connection failed !")
		} else {
			logrus.Debug("Vault connection succeeded !")
			return true
		}
	}

	logrus.Debug("Unable to connect to Vaults !")
	return false
}
