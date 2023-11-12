package proxy

import (
	"log"
	"os"
	"testing"
	"lab1/server"
	"net/http"
	"io"
	"net/url"
	"reflect"
)

type testReq struct {
	reqType     string
	path        string
	want        string
	contentType string
	body        string
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	cleanup()
}

func setup() {
	log.Println("Setup: creating server and proxy")
	server, _ := server.CreateServer("0.0.0.0", 8080, 10)
	proxy, err := CreateProxy(8081)
	if err != nil {
		panic(err)
	}

	proxy.Listen()
	server.Listen()


	go server.Serve()
	go proxy.Serve()

 
}

func cleanup() {
	 os.RemoveAll("C:/Users/Danie/Documents/Github/tda596-http/fs") 
}

func TestGetNotFound(t *testing.T) {
	tests := []testReq{
		{reqType: "GET", path: "", want: "404 Not Found"},
		{reqType: "GET", path: "/path", want: "404 Not Found"},
	}

	for _, tr := range tests {
		sendGetReq(t, tr)
	}
}

func TestPostForbidden(t *testing.T) {
	tests := []testReq{
		{reqType: "POST", path: "/path.exe", want: "403 Forbidden", body: "Hello world"},
		{reqType: "POST", path: "/path.css", want: "403 Forbidden", body: "Hello world"},
		{reqType: "POST", path: "/path.jpg", want: "403 Forbidden", body: "Hello world"},
		{reqType: "POST", path: "/path.jpeg", want: "403 Forbidden", body: "Hello world"},
		{reqType: "POST", path: "/path.html", want: "403 Forbidden", body: "Hello world"},
		{reqType: "POST", path: "/path.txt", want: "403 Forbidden", body: "Hello world"},
	}

	for _, tr := range tests {
		sendPostReq(t, tr)
	}
}


func sendGetReq(t *testing.T, tr testReq) *http.Response {
	if tr.reqType == "GET" {
		serverURL := "0.0.0.0:8080" + tr.path
    	proxyURL := "0.0.0.0:8081"

		proxy, err := url.Parse("http://" + proxyURL)
		if err != nil {
			t.Fatalf("Error parsing proxy URL: %v", err)
		}

		// Create an HTTP transport that uses the proxy
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}

		// Create an HTTP client with the transport
		client := &http.Client{
			Transport: transport,
		}

		// Make the GET request
		res, err := client.Get("http://" + serverURL)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		if !reflect.DeepEqual(string(body), tr.want) {
			t.Fatalf("\ngot:\t %s\nwant:\t %s", body, tr.want)
		}

		return res
	}
	return nil
}


func sendPostReq(t *testing.T, tr testReq) *http.Response {
	if tr.reqType == "POST" {
		serverURL := "0.0.0.0:8080" + tr.path
    	proxyURL := "0.0.0.0:8081"

		proxy, err := url.Parse("http://" + proxyURL)
		if err != nil {
			t.Fatalf("Error parsing proxy URL: %v", err)
		}

		// Create an HTTP transport that uses the proxy
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}

		// Create an HTTP client with the transport
		client := &http.Client{
			Transport: transport,
		}

		// Make the GET request
		res, err := client.Get("http://" + serverURL)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		if !reflect.DeepEqual(string(body), tr.want) {
			t.Fatalf("\ngot:\t %s\nwant:\t %s", body, tr.want)
		}

		return res
	}
	return nil
}

