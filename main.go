package main

import (
	"fmt"
	"main/core"
	"main/core/routers"
)

func main() {
	fmt.Println("Hello, world")
	server := core.Server{}
	server.LoadConfig("config.json")
	server.Create()

	server.Routers = []core.Router{
		&routers.UserRouter{Name: "users"},
	}

	server.ConnectRouters()

	server.Start(":8080")
	defer server.Dispose()
}
