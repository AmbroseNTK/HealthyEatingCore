package core

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

type Server struct {
	Echo     *echo.Echo
	Config   *Configuration
	Firebase *firebase.App
	DBClient *mongo.Client
	DB       *mongo.Database
	Routers  []Router
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
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
	server.Echo = echo.New()

	server.Echo.Validator = &Validator{validator: validator.New()}

	// Connect to database

	client, err := mongo.NewClient(options.Client().ApplyURI(server.Config.ConnectionURL))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to MongoDB")

	err = client.Connect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	server.DBClient = client
	server.DB = server.DBClient.Database(server.Config.DBName)

	log.Println("Connecting to Firebase")
	firebaseApp, firebaseError := firebase.NewApp(context.Background(), nil, option.WithServiceAccountFile(server.Config.FirebaseKeyFile))
	if firebaseError != nil {
		log.Fatal("Cannot init Firebase")
	}
	server.Firebase = firebaseApp

}

func (server *Server) Start(address string) {
	server.Echo.Logger.Fatal(server.Echo.Start(server.Config.Address))
}

func (server *Server) ConnectRouters() {
	for _, router := range server.Routers {
		router.Connect(server)
	}
}

func (server *Server) Dispose() {
	server.DBClient.Disconnect(context.TODO())
}
