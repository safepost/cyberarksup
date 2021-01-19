// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// simple does nothing except block while running the service.
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/kardianos/service"
)

//var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	// Do work here
	config := initialize()
	healthStatus := true

	var listener net.Listener
	var err error
	var isListening bool = false

	logrus.Debug("Starting health check routine ....")

	for {
		logrus.Debug("Performing checks...")
		if debug {
			// in debug mode, we look in the file status.debug and if the content is 1 we considered health OK
			content, err := ioutil.ReadFile("status.debug")
			if err != nil {
				logrus.Fatal("Unable to open status.debug file")
			}

			logrus.Debug("Content read : " + string(content))
			healthStatus = string(content) == "1"

		} else {
			healthStatus = status(config)
		}

		// check listener status
		//err := testConnection("127.0.0.1", config.port)
		//isListening := err == nil

		if healthStatus {
			if !isListening {
				// Create listener
				logrus.Debug("Starting Listener ")
				listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.port))
				if err != nil {
					logrus.Fatal("Unable to start listener !")
				}
				isListening = true
			}
		}

		if !healthStatus {
			if isListening {
				err := listener.Close()
				if err != nil {
					logrus.Fatal("Unable to stop listening ! Shutting down service")
				}
				isListening = false
			}
		}
		time.Sleep(time.Second * 10)
	}

	//	fmt.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	<-time.After(time.Second * 4)
	return nil
}

func main() {

	fmt.Printf("Running windows service\n")

	svcConfig := &service.Config{
		Name:        "GoServiceExampleStopPause",
		DisplayName: "Go Service Example: Stop Pause",
		Description: "This is an example Go service that pauses on stop.",
	}

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
	initLogger()
	//logger, err = s.Logger(nil)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	err = s.Run()
	if err != nil {
		logrus.Error(err)
	}
}
