//go:build windows

package main

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
)

type program struct{}

func (p *program) Start(service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	config := initialize()

	var listener net.Listener
	var done chan struct{}

	log.Debug("Starting health check routine ....")

	for {
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
		time.Sleep(config.getHealthCheckDuration())
	}
}

func (p *program) Stop(service.Service) error {
	// Stop should not block. Return with a few seconds.
	<-time.After(time.Second * 4)
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "Cyberark Supervision",
		DisplayName: "Cyberark Supervision Service",
		Description: "This service listen to port 38001 when Cyberark component is up",
	}

	initLogger()

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err := service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Error(err)
	}
}
