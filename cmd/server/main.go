package main

import (
	"fmt"
	"zdora/server/server"
)

func main() {
	fmt.Println("Welcome to zdora...")
	ws := server.NewWsServer()
	ws.Start()
}
