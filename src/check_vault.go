package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"net"
	"strconv"
	"strings"
	"time"
)

func testConnection(addr string, timeout time.Duration, portOptional ...int) error {
	port := 1858
	if len(portOptional) > 0 {
		port = portOptional[0]
	}

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := d.DialContext(ctx, "tcp", addr+":"+strconv.Itoa(port))

	return err
}

// Find ADDRESS filed in the Vault.ini file and return a list of IP Address (string format)
func findVaultIPAddress(iniFilePath string) ([]string, error) {
	cfg, err := ini.InsensitiveLoad(iniFilePath)
	if err != nil {
		log.Fatal("Unable to load ini file provided in configuration file : " + iniFilePath)
	}

	vaultIPAddresses := cfg.Section("").Key("ADDRESS").String()
	return strings.Split(vaultIPAddresses, ","), nil

}

func checkVault(vaultIPs []string, timeout time.Duration) bool {
	for _, ipAddr := range vaultIPs {
		err := testConnection(ipAddr, timeout)
		if err != nil {
			log.Info("Vault connection failed : " + ipAddr)
		} else {
			return true
		}
	}
	return false
}

func getVaultsIPs(iniFile string) []string {
	var vaultIPs []string

	if iniFile == "" {
		log.Info("No file provided in configuration file, trying to find Vault.ini file")
		list, err := findVaultINIFile()

		if err != nil {
			log.Fatal("Unable to find Vault INI file")
		}

		log.Info("Using Vault.ini file : " + list[0])

		vaultIPs, _ = findVaultIPAddress(list[0])
	} else {
		vaultIPs, _ = findVaultIPAddress(iniFile)
	}

	log.Info("Using addresses : " + strings.Join(vaultIPs, ","))
	return vaultIPs
}
