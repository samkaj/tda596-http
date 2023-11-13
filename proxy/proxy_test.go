package proxy

import (
	"bytes"
	"io"
	"lab1/server"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
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
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	log.Println("Setup: creating server and proxy")
	server, _ := server.CreateServer("0.0.0.0", 6060, 10)
	proxy, err := CreateProxy(6061)
	if err != nil {
		panic(err)
	}

	proxy.Listen()
	server.Listen()

	go server.Serve()
	go proxy.Serve()
}

func cleanup() {
	os.RemoveAll(os.Getenv("FS"))
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

func TestGetExistingFile(t *testing.T) {

	path := os.Getenv("FS")

	//Create files in FS directory
	server.WriteFile(path+"/test.txt", []byte("Hello world"))
	server.WriteFile(path+"/test.html", []byte("<p>Hello world</p>"))

	tests := []testReq{
		{reqType: "GET", path: "/test.txt", want: "Hello world"},
		{reqType: "GET", path: "/test.html", want: "<p>Hello world</p>"},
	}

	for _, tr := range tests {
		sendGetReq(t, tr)
	}
}

func sendGetReq(t *testing.T, tr testReq) *http.Response {
	if tr.reqType == "GET" {
		serverURL := "0.0.0.0:6060" + tr.path
		proxyURL := "0.0.0.0:6061"

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
		serverURL := "http://0.0.0.0:6060" + tr.path
		proxyURL := "http://0.0.0.0:6061"

		proxy, err := url.Parse(proxyURL)
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

		req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer([]byte(tr.body)))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		defer res.Body.Close()

		return res
	}
	return nil
}
