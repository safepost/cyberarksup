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

//func Listener(listener net.Listener, config FinalConfig ) {
//	var err error
//	listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.port))
//	if err != nil {
//		logrus.Fatal("Unable to start listener !")
//	}
//	defer listener.Close()
//	for {
//		// Wait for a connection.
//		conn, err := listener.Accept()
//		if err != nil {
//			logrus.Fatal(err)
//		}
//
//		// Handle the connection in a new goroutine.
//		// The loop then returns to accepting, so that
//		// multiple connections may be served concurrently.
//		go func(c net.Conn) {
//			// Shut down the connection.
//			c.Close()
//		}(conn)
//	}
//}

func (p *program) run() {
	// Do work here
	config := initialize()
	healthStatus := true

	var listener net.Listener
	var err error
	var isListening = false

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

				// Listener initialization
				listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.port))
				if err != nil {
					logrus.Fatal("Unable to start listener !")
				}
				defer listener.Close()

				isListening = true
			}

			conn, err := listener.Accept()
			if err != nil {
				logrus.Fatal(err)
			}
			conn.Close()

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
		Name:        "CASmartSUP",
		DisplayName: "Cyberark Smart Supervision",
		Description: "This service listen to port 38001 when Cyberark component is up",
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
