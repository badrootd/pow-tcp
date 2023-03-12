package integration

import (
	"github.com/stretchr/testify/assert"
	proto "pow-tcp/internal"
	"pow-tcp/internal/client"
	"pow-tcp/internal/server"
	"testing"
)

const address = "127.0.0.1:8081"

func Test_Integration(t *testing.T) {
	srv := server.NewServer(address)

	_ = srv.Start()

	errCh := make(chan error)
	go srv.LoopClient(errCh)

	c := client.NewClient(address)

	_ = c.Connect()

	defer c.Close()

	resp, _ := c.ReadString()

	assert.True(t, containsValue(proto.WordOfWisdom, resp))
}

func containsValue(vmap map[int]string, value string) bool {
	for _, v := range vmap {
		if v == value {
			return true
		}
	}
	return false
}
