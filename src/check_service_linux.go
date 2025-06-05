//go:build linux

package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func checkRunningService(serviceName string) (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", serviceName)
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		log.Info("Unable to get service status for: " + serviceName)
		return false, err
	}
	return true, nil
}

func checkServices(services []string) bool {
	for _, service := range services {
		isRunning, err := checkRunningService(service)
		if err != nil {
			log.Error("Service " + service + " check failed: " + err.Error())
			return false
		}
		if !isRunning {
			log.Info("Service " + service + " is not running!")
			return false
		}
	}
	return true
}
