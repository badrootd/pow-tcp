package client

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net"
	proto "pow-tcp/internal"
	"time"
)

type Client struct {
	address string
	conn    net.Conn
	rand    *rand.Rand
}

func NewClient(address string) (Client, error) {
	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return Client{}, err
	}
	r := rand.New(rand.NewSource(seed.Int64()))
	return Client{
		address: address,
		rand:    r,
	}, nil
}

func (c *Client) Connect() error {
	dial, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	c.conn = dial

	nonce, err := c.solveChallenge(c.conn)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("Failed to solve challenge")
	}

	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)
	_, err = c.conn.Write(nonceBytes)
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

func (c *Client) solveChallenge(conn net.Conn) (uint64, error) {
	arr := make([]byte, 1024)
	_, err := conn.Read(arr)
	if err != nil {
		return 0, err
	}

	algo := arr[0]
	diff := arr[1]
	prefix := arr[2:7]

	switch algo {
	case proto.SHA256:
		return c.solveSHA256(diff, prefix)
	default:
		return 0, fmt.Errorf("Algo %d not supported", algo)
	}
}

func (c *Client) solveSHA256(diff byte, prefix []byte) (uint64, error) {
	nonce := c.rand.Uint64()
	nonceBytes := make([]byte, 8)
	timeout := time.After(time.Second * 10)
	for {
		select {
		case <-timeout:
			return 0, fmt.Errorf("timed out finding nonce")
		default:
			binary.LittleEndian.PutUint64(nonceBytes, nonce)
			hash := sha256.Sum256(append(prefix, nonceBytes...))
			if proto.HasLeadingZeros(hash, int(diff)) {
				return nonce, nil
			}
			nonce++
		}
	}
}
