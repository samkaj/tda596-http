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
		{reqType: "POST", path: "/path.exe", want: "403 Forbidden", body: "Hello world"},
		{reqType: "POST", path: "/path.css", want: "200 OK", body: "Hello world"},
		{reqType: "POST", path: "/path.jpg", want: "200 OK", body: "Hello world"},
		{reqType: "POST", path: "/path.jpeg", want: "200 OK", body: "Hello world"},
		{reqType: "POST", path: "/path.html", want: "200 OK", body: "Hello world"},
		{reqType: "POST", path: "/path", want: "200 OK", body: "Hello world"},
		{reqType: "POST", path: "/path.txt", want: "200 OK", body: "Hello world"},
	}

	for _, tr := range tests {
		sendPostReq(t, tr)
	}
}

func TestPostFile(t *testing.T) {
	tests := []testReq{
		{reqType: "GET", path: "/test.txt", want: "404 Not Found"},
		{reqType: "POST", path: "/test.txt", want: "200 OK", body: "Hello world"},
		{reqType: "GET", path: "/test.txt", want: "Hello world"},
		{reqType: "POST", path: "", want: "500 Internal Server Error"},
	}

	for _, tr := range tests {
		sendPostReq(t, tr)
	}
}

func TestPostMultiple(t *testing.T) {
	tests := []testReq{
		{reqType: "POST", path: "/test1.txt", want: "200 OK", body: "1"},
		{reqType: "POST", path: "/test2.txt", want: "200 OK", body: "2"},
		{reqType: "POST", path: "/test3.txt", want: "200 OK", body: "3"},
		{reqType: "POST", path: "/test4.txt", want: "200 OK", body: "4"},
		{reqType: "POST", path: "/test5.txt", want: "200 OK", body: "5"},
		{reqType: "GET", path: "/test1.txt", want: "1"},
		{reqType: "GET", path: "/test2.txt", want: "2"},
		{reqType: "GET", path: "/test3.txt", want: "3"},
		{reqType: "GET", path: "/test4.txt", want: "4"},
		{reqType: "GET", path: "/test5.txt", want: "5"},
	}

	for _, tr := range tests {
		sendPostReq(t, tr)
	}
}

func TestNewDirectory(t *testing.T) {
	tests := []testReq{
		{reqType: "GET", path: "/testdir/test1.txt", want: "404 Not Found"},
		{reqType: "POST", path: "/testdir/test1.txt", want: "200 OK", body: "test"},
		{reqType: "GET", path: "/testdir/test1.txt", want: "test"},
	}

	for _, tr := range tests {
		sendPostReq(t, tr)
	}
}

func sendGetReq(t *testing.T, tr testReq) *http.Response {
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

		return res
	}
	return nil
}

func sendPostReq(t *testing.T, tr testReq) *http.Response {
	if tr.reqType == "POST" {

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

		return res
	}
	return nil
}

func getBodyAsString(reader io.ReadCloser) string {
	body, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}
	return string(body)
}
