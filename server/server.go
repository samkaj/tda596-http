package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Server is a simple implementation of an HTTP/1.0 web server for serving static files.
type Server struct {
	Address  string
	Port     int
	Listener net.Listener
	Sem      chan bool
}

// CreateServer tries to create an HTTP server on the specified port and address.
func CreateServer(address string, port, maxConnections int) (*Server, error) {
	// Init fs directory if nonexistent
	CreateFsDir()

	if maxConnections < 1 || maxConnections > 10 {
		return nil, fmt.Errorf("invalid amount of maximum number of connections (1-10), got %d", maxConnections)
	}

	return &Server{
		Address: address,
		Port:    port,
		Sem:     createSemaphore(maxConnections),
	}, nil
}

// Listen establishes a socket connection and listens for incoming connections.
func (s *Server) Listen() error {
	var err error
	s.Listener, err = net.Listen("tcp", s.addr())
	if err != nil {
		log.Printf("Error starting server on %s: %v", s.addr(), err)
		return err
	}
	log.Printf("Listening for connections on %s", s.addr())
	return nil
}

// Serve accepts and handles connections in goroutines.
func (s *Server) Serve() error {
	for {
		select {
		case <-s.Sem:
			conn, err := s.Listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}
			go func() {
				err := s.HandleConnection(conn)
				if err != nil {
					log.Printf("Error handling connection: %v", err)
				}
				s.Sem <- true
			}()
		}
	}
}

// HandleConnection manages incoming HTTP requests from client connections.
func (s *Server) HandleConnection(conn net.Conn) error {
	remoteAddr := conn.RemoteAddr().String()
	defer conn.Close()
	log.Printf("Handling connection from %s", remoteAddr)

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading request from %s: %v", remoteAddr, err)
		} else {
			log.Printf("Client %s closed the connection", remoteAddr)
		}
		return err
	}

	res := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	switch req.Method {
	case http.MethodGet:
		s.HandleGet(req, res)
	case http.MethodPost:
		log.Print("Handling POST")
		s.HandlePost(req, res)
	default:
		s.HandleForbidden(res)
	}

	err = res.Write(conn)
	if err != nil {
		log.Printf("Error writing response to %s: %v", remoteAddr, err)
	}
	return err
}

// DetermineContentType checks the file extension of a request.
func DetermineContentType(req *http.Request) (string, error) {
	ext := filepath.Ext(req.URL.Path)
	switch ext {
	case ".html":
		return "text/html", nil
	case ".css":
		return "text/css", nil
	case ".gif":
		return "image/gif", nil
	case ".jpeg", ".jpg":
		return "image/jpeg", nil
	case ".txt", "":
		return "text/plain", nil
	default:
		log.Printf("Invalid content type for extension: %s", ext)
		return "", fmt.Errorf("unsupported content type: %s", ext)
	}
}

// HandleGet serves GET requests.
func (s *Server) HandleGet(req *http.Request, res *http.Response) {
	contentType, err := DetermineContentType(req)
	if err != nil {
		log.Printf("Error determining content type: %v", err)
		s.HandleForbidden(res)
		return
	}

	res.Header = make(http.Header)
	res.Header.Set("Content-Type", contentType)

	filePath := filepath.Join("/Users/samkaj/code/dist/http-lab/fs", req.URL.Path)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.HandleNotFound(res)
		} else {
			log.Printf("Error reading file %s: %v", filePath, err)
			s.HandleInternalServerError(res)
		}
		return
	}

	res.Body = io.NopCloser(strings.NewReader(string(data)))
}

// HandlePost processes POST requests.
func (s *Server) HandlePost(req *http.Request, res *http.Response) {
	filePath := filepath.Join("/Users/samkaj/code/dist/http-lab/fs", req.URL.Path)

	_, err := DetermineContentType(req)
	if err != nil {
		s.HandleForbidden(res)
		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		s.HandleForbidden(res)
		return
	}

	err = WriteFile(filePath, data)
	if err != nil {
		log.Printf("Error writing to file %s with data %s: %v", filePath, data, err)
		s.HandleInternalServerError(res)
		return
	}

	res.Body = io.NopCloser(strings.NewReader("200 OK"))
}

// HandleNotFound builds a 404 Not Found response.
func (s *Server) HandleNotFound(res *http.Response) {
	res.Status = "404 Not Found"
	res.StatusCode = 404
	res.Body = io.NopCloser(strings.NewReader("404 Not Found"))
}

// HandleForbidden handles forbidden HTTP methods by setting a 403 Forbidden response.
func (s *Server) HandleForbidden(res *http.Response) {
	res.Status = "403 Forbidden"
	res.StatusCode = 403
	res.Body = io.NopCloser(strings.NewReader("403 Forbidden"))
}

// HandleInternalServerError builds a 500 Internal Server Error response.
func (s *Server) HandleInternalServerError(res *http.Response) {
	res.Status = "500 Internal Server Error"
	res.StatusCode = 500
	res.Body = io.NopCloser(strings.NewReader("500 Internal Server Error"))
}

func CreateBody(text string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(text + "\n"))
}

// Close attempts to close the server's listener and logs any error.
func (s *Server) Close() {
	if err := s.Listener.Close(); err != nil {
		log.Printf("Error closing server listener: %v", err)
	}
}

// addr returns the server's address and port as a string.
func (s *Server) addr() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}

// createSemaphore creates a channel to control the number of active connections.
func createSemaphore(size int) chan bool {
	sem := make(chan bool, size)
	for i := 0; i < size; i++ {
		sem <- true
	}
	return sem
}
