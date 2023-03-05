package client

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net"
	proto "pow-tcp/internal"
)

type Client struct {
	address string
	conn    net.Conn
}

func NewClient(address string) Client {
	return Client{
		address: address,
	}
}

func (c *Client) Connect() error {
	dial, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	c.conn = dial

	nonce, err := solveChallenge(c.conn)
	if err != nil {
		return err
	}

	if nonce == nil {
		return fmt.Errorf("Failed to solve challenge")
	}

	_, err = c.conn.Write(nonce)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ReadString() (string, error) {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func solveChallenge(conn net.Conn) ([]byte, error) {
	arr := make([]byte, 1024)
	_, err := conn.Read(arr)
	if err != nil {
		return nil, err
	}

	algo := arr[0]
	fmt.Printf("Algo: %d\n ", algo)
	diff := arr[1]
	fmt.Printf("Difficulty: %d\n ", diff)
	prefix := arr[2:7]
	fmt.Println("prefix: " + string(prefix))

	var nonce []byte
	switch algo {
	case proto.SHA256:
		nonce = solveSHA256(diff, prefix)
	default:
		return nil, fmt.Errorf("Algo %d not supported", algo)
	}

	fmt.Println("nonce: " + string(nonce))
	return nonce, nil
}

func solveSHA256(diff byte, prefix []byte) []byte {
	nonce := proto.RandSeq(5)
	var count int
	for {
		count++
		hash := sha256.Sum256(append(prefix, nonce...))
		if proto.HasLeadingZeros(hash, int(diff)) {
			// TODO: remove
			fmt.Printf("Hash: %s\n", hash)
			fmt.Printf("Solved in %d attempts\n", count)
			return nonce
		}
		_, err := rand.Read(nonce)
		if err != nil {
			fmt.Println("err")
		}
	}
}
