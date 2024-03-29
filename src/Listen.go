// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func accept(l net.Listener, done chan struct{}) {
	for {
		conn, err := l.Accept()
		select {
		case <-done:
			return
		default:
		}
		if err != nil {
			panic(err)
		}
		go func(c net.Conn) {
			_ = c.Close()
		}(conn)
	}
}

func healthcheck(config FinalConfig) bool {
	if debug {
		// in debug mode, we look in the file status.debug and if the content is 1 we considered health OK
		content, err := ioutil.ReadFile("status.debug")
		if err != nil {
			log.Fatal("Unable to open status.debug file")
		}

		log.Debug("Content read : " + string(content))
		return string(content) == "1"

	} else {
		return status(config)
	}
}

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

func main() {
	svcConfig := &service.Config{
		Name:        "CASmartSUP",
		DisplayName: "Cyberark Smart Supervision",
		Description: "This service listen to port 38001 when Cyberark component is up",
	}

	initLogger()

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
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
