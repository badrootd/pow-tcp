package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"pow-tcp/internal/client"
)

func main() {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		log.Errorf("Failed to get host env var:")
		return
	}

	log.Info("Connecting to %s", host)
	c, err := client.NewClient(host)
	if err != nil {
		log.Errorf("Error creating client:", err)
		return
	}

	if err = c.Connect(); err != nil {
		log.Errorf("Error connecting to server:", err)
		return
	}

	defer c.Close()

	resp, _ := c.ReadString()
	fmt.Println(resp)
}
