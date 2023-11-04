package main

import (
	"fmt"
	"lab1/server"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
	}

	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	server, err := server.CreateServer(host, port, 3)
	if err != nil {
		fmt.Printf("failed to start server with error: %v", err)
		os.Exit(1)
	}

	server.Listen()
	server.Serve()
	defer server.Close()
}

func printUsage() {
	fmt.Println("Usage: http_server <host> <port>")
	os.Exit(1)
}
