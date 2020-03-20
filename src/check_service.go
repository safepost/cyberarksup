package main

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func checkRunningService(serviceName string, manager *mgr.Mgr) (bool, error) {
	service, err := manager.OpenService(serviceName)
	if err != nil {
		return false, err
	}

	status, err := service.Query()
	if err != nil {
		return false, err
	}

	return status.State == svc.Running, nil
}

func checkServices(services []string) bool {

	// Check running services
	manager, err := mgr.Connect()
	if err != nil {
		logrus.Info("Unable to connect to service manager !")
		return false
		//panic(err)
	}
	defer manager.Disconnect()

	for _, service := range services {
		isRunning, err := checkRunningService(service, manager)
		if err != nil {
			logrus.Fatal("Service " + service + " does not exists or name is invalid ! Exiting.")
		}
		if !isRunning {
			logrus.Info("Service " + service + " is not running! ")
			return false
		}
	}

	return true
}
