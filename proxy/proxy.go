package proxy

import (
	"bufio"
	"fmt"
	"lab1/server"
	"log"
	"net"
	"net/http"
	"strings"
)

// Proxy is a wrapper for a regular server with additional
//   - restraints - only GET:s are allowed,
//   - functionality - acts on behalf of the client by making the 
//     requests and passing back the response.
type Proxy struct {
	proxyServer *server.Server
}

// Creates a proxy on the given port, listening on any address.
// Returns any errors that occurred.
func CreateProxy(port int) (*Proxy, error) {
	server, err := server.CreateServer("0.0.0.0", port, 10)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		proxyServer: server,
	}, nil
}

// Wrapper for server implementation.
// See server/server.go.
func (p *Proxy) Listen() error {
	err := p.proxyServer.Listen()
	if err != nil {
		return err
	}

	return nil
}

// Wrapper for server implementation.
// See server/server.go.
func (p *Proxy) Serve() error {
	for {
		if <-p.proxyServer.Sem {
			conn, err := p.proxyServer.Listener.Accept()
			if err != nil {
				return err
			}

			go p.HandleConnection(conn)
		}
	}
}

func (p *Proxy) HandleConnection(conn net.Conn) error {
	defer conn.Close()
	log.Printf("handling connection via proxy from %s\n", conn.RemoteAddr().String())

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return err
	}

	// Only allow HTTP GET.
	if req.Method != http.MethodGet {
        p.SendNotImplemented(conn)
		return fmt.Errorf("received forbidden HTTP method: %s", req.Method)
	}

	// Act on behalf of the client (proxy user).
    res, err := p.SendRequestToServer(req)
	if err != nil {
		return err
	}

	// Send back the response to the proxy user.
	p.SendResponseToClient(conn, res)
	p.proxyServer.Sem <- true
	return nil
}

// Sends a HTTP GET request to the server and returns it and any
// errors that occured.
func (p *Proxy) SendRequestToServer(req *http.Request) (*http.Response, error) {
	path := req.URL.Path
	if strings.HasSuffix(path, "/") {
		path = strings.TrimRight(path, "/")
	}

	res, err := http.Get(req.RequestURI)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return res, nil
}

// Sends a 501 - Not Implemented to the client.
func (p *Proxy) SendNotImplemented(conn net.Conn) {
	res := &http.Response{
		Status:     "501 Not Implemented",
		StatusCode: 501,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
	}
    res.Write(conn)
}

// Sends back the response acquired from the server to the client
// using the proxy.
func (p *Proxy) SendResponseToClient(conn net.Conn, res *http.Response) error {
	err := res.Write(conn)
	if err != nil {
		return err
	}
	return nil
}

// Wrapper for closing the server.
func (p *Proxy) Close() {
	p.proxyServer.Close()
}
