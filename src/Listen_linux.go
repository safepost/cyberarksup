//go:build linux
// +build linux

package main

import (
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"time"
)

type program struct{}

func (p *program) run() {
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
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	<-time.After(time.Second * 4)
	return nil
}
