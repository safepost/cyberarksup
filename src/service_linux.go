//go:build linux
// +build linux

package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/takama/daemon"

	//	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
)

const (

	// name of the service
	name        = "Cyberark Supervision"
	description = "Cyberark Monitoring service"

	// port which daemon should be listen
	port = ":9977"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Do something, call your goroutines, etc
	config := initialize()

	var listener net.Listener
	var done chan struct{}

	log.Debug("Starting health check routine ....")

	for {
		log.Debug("Performing checks...")

		if healthcheck(config) {
			if listener == nil {
				var err error
				listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.port))
				if err != nil {
					panic(err)
				}
				done = make(chan struct{})
				go accept(listener, done)
			}
		} else {
			if listener != nil {
				close(done)
				err := listener.Close()
				if err != nil {
					panic(err)
				}
				done = nil
				listener = nil
			}
		}
		time.Sleep(time.Second * 10)
	}

	// never happen, but need to complete code
	// return usage, nil
}

// Accept a client connection and collect it in a channel

func main() {
	srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
	if err != nil {
		log.Error("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		log.Error(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
