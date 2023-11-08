package main

import (
	"lab1/proxy"
	"log"
	"os"
	"strconv"
)

func main() {
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
