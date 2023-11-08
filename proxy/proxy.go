package proxy

import (
	"bufio"
	"fmt"
	"io"
	"lab1/server"
	"log"
	"net"
	"net/http"
	"strings"
)

type Proxy struct {
	proxyServer   *server.Server
}

func CreateProxy(port int) (*Proxy, error) {
    server, err := server.CreateServer("0.0.0.0", port, 10)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		proxyServer: server,
	}, nil
}

func (p *Proxy) Listen() error {
	err := p.proxyServer.Listen()
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) Serve() error {
	for {
		if <-p.proxyServer.Sem {
			conn, err := p.proxyServer.Listener.Accept()
			if err != nil {
				return err
			}

			go p.HandleGetConnection(conn)
		}
	}
}

func (p *Proxy) HandleGetConnection(conn net.Conn) error {
	log.Printf("handling connection via proxy from %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return err
	}

	body := io.NopCloser(strings.NewReader(""))
	res := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       body,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
	}

	if req.Method != http.MethodGet {
		p.proxyServer.HandleForbidden(res)
		return fmt.Errorf("received forbidden HTTP method: %s", req.Method)
	}

    res, err = p.SendGetToServer(req)
	if err != nil {
		return err
	}

	p.SendResponseToClient(conn, res)
	p.proxyServer.Sem <- true
	return nil
}

// TODO: better naming
func (p *Proxy) SendGetToServer(req *http.Request) (*http.Response, error) {
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

func (p *Proxy) SendResponseToClient(conn net.Conn, res *http.Response) error {
	err := res.Write(conn)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) Close() {
	p.proxyServer.Close()
}
