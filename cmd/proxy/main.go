package main

import (
	"lab1/proxy"
	"log"
)

func main() {
    log.SetPrefix("[PROXY] ")
	proxy, err := proxy.CreateProxy("localhost", "localhost", 8080, 6060)
	if err != nil {
		panic(err)
	}

	proxy.Listen()
	proxy.Serve()
	defer proxy.Close()
}
