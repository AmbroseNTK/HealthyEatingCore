package main

import (
	"fmt"
	"main/core"
)

func main() {
	fmt.Println("Hello, world")
	server := core.Server{}
	server.LoadConfig("config.json")
	server.Create()
	server.Start(":8080")
	defer server.Dispose()
}
