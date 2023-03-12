package main

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"pow-tcp/internal/server"
	"time"
)

const address = ":8081"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	srv := server.NewServer(address)

	err := srv.Start()
	if err != nil {
		log.Error("Shutting down, failed to start server: %s", err)
		return
	}
	defer srv.Close()
	log.Info("Listening on %s", address)

	errCh := make(chan error)

	go handleErrors(errCh)

	srv.LoopClient(errCh)
}

func handleErrors(errCh <-chan error) {
	for {
		select {
		case err := <-errCh:
			log.Error(err)
		}
	}
}
