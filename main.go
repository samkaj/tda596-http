package main

import (
	"fmt"
	"lab1/server"
)

func main() {
	fmt.Println("Hello world")
	server, err := server.CreateServer("localhost", 8080, 3)
	if err != nil {
		panic(err)
	}

	server.Listen()
	defer server.Close()
}
