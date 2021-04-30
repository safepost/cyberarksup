package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func checkRunningService(serviceName string) (bool, error) {
	cmd := exec.Command("systemctl", "check", serviceName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		} else {
			log.Info("Unable to get psmp status")
			return false, err
		}
	}
	return true, nil
}

func checkServices(services []string) bool {
	// Check running services
	for _, service := range services {
		isRunning, err := checkRunningService(service)
		if err != nil {
			log.Panic("Service " + service + " does not exists or name is invalid ! Exiting.")
			panic(err)
		}
		if !isRunning {
			log.Info("Service " + service + " is not running! ")
			return false
		}
	}
	return true
}
