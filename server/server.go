package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strings"
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

	body := io.NopCloser(strings.NewReader("hello"))

	res := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       body,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
	}

	switch method := req.Method; method {
	case http.MethodGet:
		s.HandleGet(req, res)
	case http.MethodPut:
		s.HandlePut(req, res)
	default:
		s.HandleBadRequest(req, res)
	}

	res.Write(conn)

	return nil
}

// Checks the file extension of a request and returns the Content-Type.
// It returns any error that occured.
// Allowed Content-Types: text/html, text/plain, image/gif, image/jpeg, image/jpeg, or text/css.
func (s *Server) DetermineContentType(req http.Request) (string, error) {

	switch ct := path.Ext(req.URL.Path); ct {
	case ".html":
		return "text/html", nil
	case "":
		return "text/plain", nil
	case ".css":
		return "text/css", nil
	case ".gif":
		return "image/gif", nil
	case ".jpeg":
		return "image/jpeg", nil
	case ".jpg":
		return "image/jpg", nil
	}
	return "", fmt.Errorf("invalid content type")
}

// Handles HTTP GET requests.
func (s *Server) HandleGet(req *http.Request, res *http.Response) {
	fmt.Println("Got GET")

}

// Handles HTTP PUT requests.
func (s *Server) HandlePut(req *http.Request, res *http.Response) {
	fmt.Println("Got PUT")
}

// Handles forbidden HTTP methods.
func (s *Server) HandleBadRequest(req *http.Request, res *http.Response) {

	// Bad request
	res.Status = "400 Bad Request"
	res.StatusCode = 400
	res.Body = CreateBody("400 Bad Request")
}

func CreateBody(text string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(text + "\n"))
}

// Closes the server's listener.
func (s *Server) Close() {
	s.Listener.Close()
}

// Returns the server's address and port as a string.
func (s *Server) addr() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}
