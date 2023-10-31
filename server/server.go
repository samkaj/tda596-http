package server

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	Address        string
	Port           int
	MaxConnections int
	Listener       net.Listener
}

// Tries to create an HTTP server on the specific port and address that allows
// at most maxConnections concurrent connections (at most 10).
// It returns the server and any error that occured.
func CreateServer(address string, port, maxConnections int) (*Server, error) {
	if maxConnections < 1 || maxConnections > 10 {
		return nil, fmt.Errorf("invalid amount of maximum number of connections (1-10), got %d", maxConnections)
	}

	return &Server{
		Address:        address,
		Port:           port,
		MaxConnections: maxConnections,
	}, nil
}

// Establishes a socket connection and listens for connections.
// It returns any error that occured.
func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.addr())
	if err != nil {
		return err
	}
	s.Listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.HandleConnection(conn)
	}
}

// Given an established connection to a client, this method will handles incoming HTTP requests from the client
func (s *Server) HandleConnection(conn net.Conn) error {
	// Defer is run after the HandleConnection routine is done
	defer conn.Close()

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return err
	}

	switch method := req.Method; method {
	case http.MethodGet:
		s.HandleGet(req)
	case http.MethodPut:
		s.HandlePut(req)
	default:
		fmt.Println("FIXME: bad request...")
	}

	return nil
}

// Handles HTTP GET requests.
func (s *Server) HandleGet(req *http.Request) {

}

// Handles HTTP PUT requests.
func (s *Server) HandlePut(req *http.Request) {

}

// Handles forbidden HTTP methods.
func (s *Server) HandleBadRequest(req *http.Request) {
}

// Closes the server's listener.
func (s *Server) Close() {
	s.Listener.Close()
}

// Returns the server's address and port as a string.
func (s *Server) addr() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}
