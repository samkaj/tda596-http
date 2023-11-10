package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
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
	log.Println("Setup: creating server")
	server, _ := CreateServer("0.0.0.0", 8080, 10)
	server.Listen()
	go server.Serve()
}

func cleanup() {
	os.RemoveAll("/Users/samkaj/code/dist/http-lab/fs")
}

func TestGetNotFound(t *testing.T) {
	tests := []testReq{
		{reqType: "GET", path: "", want: "404 Not Found"},
		{reqType: "GET", path: "/path", want: "404 Not Found"},
		{reqType: "GET", path: "/hej", want: "404 Not Found"},
	}

	for _, tr := range tests {
		sendGetReq(t, tr)
	}
}

func TestPostContentType(t *testing.T) {
	tests := []testReq{
		{reqType: "POST", path: "/path.tsdfasdfxt", want: "403 Forbidden", body: "Hello world"},
		{reqType: "POST", path: "/path.css", want: "200 OK", body: "Hello world"},
	}

	for _, tr := range tests {
		sendPostReq(t, tr)
	}
}

func sendGetReq(t *testing.T, tr testReq) {
	if tr.reqType == "GET" {
		res, err := http.Get("http://0.0.0.0:8080" + tr.path)
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
	}
}

func sendPostReq(t *testing.T, tr testReq) {
	if tr.reqType == "POST" {
		if tr.body == "" {
			t.Fatal("No body provided in post request")
		}

		inbody := strings.NewReader(tr.body)
		res, err := http.Post("http://0.0.0.0:8080"+tr.path, tr.contentType, inbody)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		if !reflect.DeepEqual(string(body), tr.want) {
			t.Fatalf("\ngot:\t %s\nwant:\t %s", body, tr.want)
		}
	}
}
