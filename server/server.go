package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	Address  string
	Port     int
	Listener net.Listener
	sem      chan bool
}

// Tries to create an HTTP server on the specific port and address that allows
// at most maxConnections concurrent connections (at most 10).
// It returns the server and any error that occured.
func CreateServer(address string, port, maxConnections int) (*Server, error) {
	if maxConnections < 1 || maxConnections > 10 {
		return nil, fmt.Errorf("invalid amount of maximum number of connections (1-10), got %d", maxConnections)
	}

	return &Server{
		Address: address,
		Port:    port,
		sem:     createSemaphore(maxConnections),
	}, nil
}

// Establishes a socket connection and listens for connections.
// It returns any error that occured.
func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.addr())
	if err != nil {
		return err
	}

	log.Printf("listening for connections on %s", s.addr())
	s.Listener = listener
	return nil
}

// Accepts and handles connections in goroutines.
func (s *Server) Serve() error {
	for {
		if <-s.sem {
			conn, err := s.Listener.Accept()
			if err != nil {
				return err
			}
			go s.HandleConnection(conn)
		}
	}
}

// Handles incoming HTTP requests from client connection.
func (s *Server) HandleConnection(conn net.Conn) error {
	log.Printf("handling connection from %s\n", conn.LocalAddr().String())
	time.Sleep(time.Millisecond * 500)

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
	case http.MethodPost:
		s.HandlePost(req, res)
	default:
		s.HandleBadRequest(req, res)
	}

	res.Write(conn)
	s.sem <- true
	return nil
}

// Checks the file extension of a request and returns the Content-Type.
// It returns any error that occured.
// Allowed Content-Types: text/html, text/plain, image/gif, image/jpeg, image/jpeg, or text/css.
func DetermineContentType(req *http.Request) (string, error) {
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
	contentType, err := DetermineContentType(req)
	if err != nil {
		s.HandleBadRequest(req, res)
		return
	}

	// Set the header content type based on the extension that has been requested
	res.Header = make(http.Header)
	res.Header.Set("Content-Type", contentType)

	// Put together the path to the file with path.Join
	path := filepath.Join("fs", req.URL.Path)
	if err != nil {
		s.HandleNotFound(res)
		return
	}
	data, err := GetFile(path)
	if err != nil {
		s.HandleNotFound(res)
		return
	}

	res.Body = CreateBody(string(data))
}

// Handles HTTP POST requests.
func (s *Server) HandlePost(req *http.Request, res *http.Response) {
	// Get filename and path
	filename := path.Base(req.URL.Path)
	path := filepath.Join("fs", filename)

	// Get file data
	data, err := io.ReadAll(req.Body)
	if err != nil {
		s.HandleBadRequest(req, res)
	}
	// Write to fs
	WriteFile(path, data)

	// Build response
	res.Status = "200 OK"
	res.StatusCode = 200
	res.Body = CreateBody("200 OK")
}

// Builds a 404 Not Found response.
func (s *Server) HandleNotFound(res *http.Response) {
	res.Status = "404 Not Found"
	res.StatusCode = 404
	res.Body = CreateBody("404 Not Found")
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

// Creates and initializes channel for controlling
// number of active connections.
func createSemaphore(size int) chan bool {
	sem := make(chan bool, size)
	for i := 0; i < size; i++ {
		sem <- true
	}
	return sem
}
