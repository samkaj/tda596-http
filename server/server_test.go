package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

type testReq struct {
	name 	  	string
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
	log.Println("Setup: creating server")
	server, _ := CreateServer("0.0.0.0", 8080, 10)
	server.Listen()
	go server.Serve()
}

func cleanup() {
	fsPath := os.Getenv("FS")
	os.RemoveAll(fsPath)
}

func TestGetNotFound(t *testing.T) {
	tests := []testReq{
		{name: "Get non existant file",reqType: "GET", path: "", want: "404 Not Found"},
		{name: "Get non existant file",reqType: "GET", path: "/path", want: "404 Not Found"},
		{name: "Get non existant file",reqType: "GET", path: "/hej", want: "404 Not Found"},
	}

	for _, tr := range tests {
		sendReq(t, tr)
	}
}

func TestPostContentType(t *testing.T) {
	tests := []testReq{
		{name: "Post .exe", reqType: "POST", path: "/path.exe", want: "400 Bad Request", body: "Hello world"},
		{name: "Post .css",reqType: "POST", path: "/path.css", want: "200 OK", body: "Hello world"},
		{name: "Post .jpg",reqType: "POST", path: "/path.jpg", want: "200 OK", body: "Hello world"},
		{name: "Post .jpeg",reqType: "POST", path: "/path.jpeg", want: "200 OK", body: "Hello world"},
		{name: "Post .html",reqType: "POST", path: "/path.html", want: "200 OK", body: "Hello world"},
		{name: "Post plain",reqType: "POST", path: "/path", want: "200 OK", body: "Hello world"},
		{name: "Post .txt",reqType: "POST", path: "/path.txt", want: "200 OK", body: "Hello world"},
	}

	for _, tr := range tests {
		sendReq(t, tr)
	}
}

func TestPostFile(t *testing.T) {
	tests := []testReq{
		{name: "Get non existant file", reqType: "GET", path: "/test.txt", want: "404 Not Found"},
		{name: "Post test.txt",reqType: "POST", path: "/test.txt", want: "200 OK", body: "Hello world"},
		{name: "Get test.txt",reqType: "GET", path: "/test.txt", want: "Hello world"},
		{name: "Post without specifying filename", reqType: "POST", path: "", want: "400 Bad Request"},
	}

	for _, tr := range tests {
		sendReq(t, tr)
	}
}

func TestPostMultiple(t *testing.T) {
	tests := []testReq{
		{name: "Post test1.txt",reqType: "POST", path: "/test1.txt", want: "200 OK", body: "1"},
		{name: "Post test2.txt",reqType: "POST", path: "/test2.txt", want: "200 OK", body: "2"},
		{name: "Post test3.txt",reqType: "POST", path: "/test3.txt", want: "200 OK", body: "3"},
		{name: "Post test4.txt",reqType: "POST", path: "/test4.txt", want: "200 OK", body: "4"},
		{name: "Post test5.txt",reqType: "POST", path: "/test5.txt", want: "200 OK", body: "5"},
		{name: "Get test1.txt",reqType: "GET", path: "/test1.txt", want: "1"},
		{name: "Get test2.txt",reqType: "GET", path: "/test2.txt", want: "2"},
		{name: "Get test3.txt",reqType: "GET", path: "/test3.txt", want: "3"},
		{name: "Get test4.txt",reqType: "GET", path: "/test4.txt", want: "4"},
		{name: "Get test5.txt",reqType: "GET", path: "/test5.txt", want: "5"},
	}

	for _, tr := range tests {
		sendReq(t, tr)
	}
}

func TestNewDirectory(t *testing.T) {
	tests := []testReq{
		{name: "Get non existant file from non existant directory",reqType: "GET", path: "/testdir/test1.txt", want: "404 Not Found"},
		{name: "Post test1.txt file",reqType: "POST", path: "/testdir/test1.txt", want: "200 OK", body: "test"},
		{name: "Get test1.txt file",reqType: "GET", path: "/testdir/test1.txt", want: "test"},
	}

	for _, tr := range tests {
		sendReq(t, tr)
	}
}


func TestNotImplemented(t *testing.T) {
	tests := []testReq{
		{name: "Send DELETE",reqType: "DELETE", path: "/testdir/test1.txt", want: "501 Not Implemented"},
		{name: "Send PUT",reqType: "PUT", path: "/testdir/test1.txt", want: "501 Not Implemented"},
		{name: "Send HEAD",reqType: "HEAD", path: "/testdir/test1.txt", want: "501"},
	}

	for _, tr := range tests {
		sendReq(t, tr)
	}
}

func sendReq(t *testing.T, tr testReq) *http.Response {
	if tr.reqType == http.MethodPost {

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
	}else if tr.reqType == http.MethodGet {
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
			failTest(t, tr, string(body))
		}
	
		return res
	} else if tr.reqType == http.MethodDelete {
		req, err := http.NewRequest(http.MethodDelete,"http://0.0.0.0:8080" + tr.path, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		client := &http.Client{}

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
	
		defer res.Body.Close()
	
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
	
		if !reflect.DeepEqual(string(body), tr.want) {
			failTest(t, tr, string(body))
		}
	
		return res

	} else if tr.reqType == http.MethodPut {
		inbody := strings.NewReader(tr.body)
		req, err := http.NewRequest(http.MethodPut,"http://0.0.0.0:8080" + tr.path, inbody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		client := &http.Client{}

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
	
		defer res.Body.Close()
	
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
	
		if !reflect.DeepEqual(string(body), tr.want) {
			failTest(t, tr, string(body))
		}
	
		return res

	} else if tr.reqType == http.MethodHead {
		res, err := http.Head("http://0.0.0.0:8080" + tr.path)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
	
		defer res.Body.Close()
	
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		
		code, err := strconv.Atoi(tr.want)
		if !reflect.DeepEqual(res.StatusCode, code) {
			failTest(t, tr, string(body))
		}

		return res
	}
	
	t.Fatalf("Non existant request method: %v", tr.reqType)
	return nil
}


func failTest(t *testing.T, tr testReq, got string) {
	
	t.Fatalf("\nName: %s \ngot:\t %s\nwant:\t %s",tr.name, got, tr.want)
}

func getBodyAsString(reader io.ReadCloser) string {
	body, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}
	return string(body)
}
