package main

import (
	"fmt"
	"lab1/proxy"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("failed to load environment variables, create the .env file")
	}

	log.SetPrefix("[PROXY] ")
	port, err := strconv.Atoi(os.Args[1])
	proxy, err := proxy.CreateProxy(port)
	if err != nil {
		panic(err)
	}

	proxy.Listen()
	proxy.Serve()
	defer proxy.Close()
}
