package server

import (
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
)

var server *Server

type testReq struct {
	reqType string
	path string
	want  string
	Body  string
}

func setup() *Server{
	server, _ = CreateServer("localhost", 8080, 10)
	server.Listen()
	go server.Serve()
	return server;

}

func shutdown() {
	server.Close()
}

func TestMain(m *testing.M) {
    setup()
    code := m.Run() 
    shutdown()
    os.Exit(code)
}

func TestListen(t *testing.T){

	if server.Listener == nil {
		t.Errorf("server.Listener is nil")
	}
}


func TestGetNotFound(t *testing.T) {
	tests := []testReq{
        {reqType: "GET", path: "", want: "HTTP/1.0 404 Not Found\n"},
        {reqType: "GET", path: "/path", want: "HTTP/1.0 404 Not Found\n"},
        {reqType: "GET", path: "/hej", want: "HTTP/1.0 404 Not Found\n"},
    }

	for _, tr := range tests {
        sendReq(t, tr)
    }
}

func TestPost(t *testing.T){
	tests := []testReq{
        {reqType: "POST", path: "/test.txt", want: "HTTP/1.0 404 Not Found\n", Body: "Hello World"},
    }

	for _, tr := range tests {
        sendReq(t, tr)
    }

}



func sendReq(t *testing.T, tr testReq ) {
    // Create a pair of connected network connections
    conn1, conn2 := net.Pipe()

    // Run HandleConnection in a goroutine
    go func() {
        err := server.HandleConnection(conn1)
        if err != nil {
            t.Errorf("HandleConnection returned error: %v", err)
        }
    }()
    // Use conn2 to send and receive data
	var req *http.Request
	var err error
	if(tr.reqType == "POST"){
		req, err = http.NewRequest(tr.reqType, tr.path, strings.NewReader(tr.Body))
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}
	}

	if tr.reqType == "GET" {
		req, err = http.NewRequest(tr.reqType, "/path", nil)
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}
	}
   
	err = req.Write(conn2)
	if err != nil {
		t.Fatalf("Could not write HTTP request: %v", err)
	}

    // Read the response
    buf := make([]byte, 1024)
    n, err := conn2.Read(buf)
    if err != nil {
        t.Fatalf("Failed to read from connection: %v", err)
    }
    // Check the response
	
    if string(buf[:n]) != tr.want {
        t.Errorf("Unexpected response: %s", string(buf[:n]))
		t.Errorf("Expected response: %s", tr.want)
    }
}
