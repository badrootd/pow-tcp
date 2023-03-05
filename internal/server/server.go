package server

import (
	"crypto/sha256"
	_ "crypto/sha256"
	"fmt"
	"math/rand"
	"net"
	proto "pow-tcp/internal"
)

type Server struct {
	address  string
	listener net.Listener
}

func NewServer(address string) *Server {
	return &Server{address: address}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("Error listening: %v", err)
	}

	s.listener = listener

	return nil
}

func (s *Server) LoopClient(errch chan<- error) {
	for {
		s.nextClient(errch)
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) nextClient(errch chan<- error) {
	conn, err := s.listener.Accept()
	if err != nil {
		errch <- fmt.Errorf("Error accepting: %w", err)
		return
	}

	go handleConnection(conn, errch)
}

func handleConnection(conn net.Conn, errch chan<- error) {
	defer conn.Close()

	prefix, err := challenge(conn)
	if err != nil {
		errch <- err
		return
	}

	var ok bool
	ok, err = verify(conn, prefix)
	if err != nil {
		errch <- err
		return
	}

	if !ok {
		return
	}

	n := rand.Intn(3)
	word, _ := proto.WordOfWisdom[n]
	_, err = conn.Write([]byte(word))
}

func challenge(conn net.Conn) ([]byte, error) {
	arr := make([]byte, 2)
	arr[0] = proto.SHA256
	arr[1] = proto.Difficulty

	prefix := proto.RandSeq(5)
	arr = append(arr, prefix...)

	_, err := conn.Write(arr)

	return prefix, err
}

func verify(conn net.Conn, prefix []byte) (bool, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return false, err
	}

	nonce := string(buf[:n])

	hash := sha256.Sum256(append(prefix, nonce...))
	if proto.HasLeadingZeros(hash, proto.Difficulty) {
		return true, nil
	}

	return false, nil
}
