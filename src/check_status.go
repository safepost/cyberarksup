package main

import "github.com/sirupsen/logrus"

func status(config FinalConfig) bool {

	logrus.Debug("Starting health tests...")

	for _, disk := range config.disks {
		diskStatus := DiskUsage(disk)
		logrus.Debug(diskStatus)
		if !diskStatus {
			logrus.Info("Disk " + disk + " available space is less than 10% !")
			return false
		}
	}

	vaultConn := checkVault(config.vaultIPs)
	if !vaultConn {
		logrus.Info("Connexion to Vaults failed !")
		return false
	}

	serviceStatus := checkServices(config.services)
	if !serviceStatus {
		logrus.Info("Service check failed !")
		return false
	}

	logrus.Debug("All checks went well !")
	return true

}
