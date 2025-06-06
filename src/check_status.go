package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func status(config FinalConfig) bool {
	timeout := config.getVaultTimeOutDuration()

	// Check Vault
	vaultConn := checkVault(config.vaultIPs, timeout)
	if !vaultConn {
		log.Info("[FAIL] Connexion to Vaults failed !")
		return false
	} else {
		log.Debug(fmt.Sprintf("[PASS] Connexion to Vault succeeded, IP [%s]", config.vaultIPs))
	}

	// Check Services
	serviceStatus := checkServices(config.services)
	if !serviceStatus {
		log.Info("[FAIL] Service check failed !")
		return false
	} else {
		log.Debug("[PASS] Configured services were running")
	}

	// Check Disks
	for _, disk := range config.disks {
		diskStatus := DiskUsage(disk)
		if !diskStatus {
			log.Info("Disk " + disk + " available space is less than 10% !")
			return false
		} else {
			log.Debug("[PASS] Disk usage check passed")
		}
	}

	log.Debug("[SUCCESS] All checks went well.")
	return true

}
