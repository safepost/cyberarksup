package main

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func healthcheck(config FinalConfig) bool {
	if debug {
		// in debug mode, we look in the file status.debug and if the content is 1 we considered health OK
		content, err := os.ReadFile("status.debug")
		if err != nil {
			log.Fatal("Unable to open status.debug file")
		}

		log.Debug("Content read : " + string(content))
		return string(content) == "1"

	} else {
		return status(config)
	}
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
