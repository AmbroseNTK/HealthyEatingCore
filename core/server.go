package core

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	e      *echo.Echo
	Config *Configuration
	db     *mongo.Client
}

func (server *Server) LoadConfig(configFile string) {
	if configFile == "" {
		server.Config = &Configuration{
			Address: ":8080",
		}
		return
	}
	server.Config = &Configuration{}
	server.Config.Load(configFile)
}

func (server *Server) Create() {
	server.e = echo.New()

	// Connect to database

	client, err := mongo.NewClient(options.Client().ApplyURI(server.Config.ConnectionURL))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	server.db = client

	server.e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello, world",
		})
	})
}

func (server *Server) Start(address string) {
	server.e.Logger.Fatal(server.e.Start(server.Config.Address))
}

func (server *Server) Dispose() {
	server.db.Disconnect(context.TODO())
}
