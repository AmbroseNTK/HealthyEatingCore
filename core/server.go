package core

import (
	"context"
	"log"
	"main/core/middlewares"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

type Server struct {
	Echo              *echo.Echo
	Config            *Configuration
	Firebase          *firebase.App
	Auth              *auth.Client
	DBClient          *mongo.Client
	DB                *mongo.Database
	Routers           []Router
	AuthMiddleware    func(next echo.HandlerFunc) echo.HandlerFunc
	AuthWiddlewareJWT *middlewares.AuthMiddleware
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

	auth, authError := firebaseApp.Auth(context.TODO())

	if authError != nil {
		log.Fatal("Cannot load auth module")
	}
	server.Auth = auth

	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idToken := c.Request().Header.Get("Authorization")
			if idToken == "" {
				idToken = c.QueryParam("token")
			}
			if idToken == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing ID Token")
			}
			token, tokenError := auth.VerifyIDToken(context.Background(), idToken)
			if tokenError != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, tokenError.Error())
			}
			c.Set("user", token)
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}

	}

	server.AuthMiddleware = authMiddleware
	server.AuthWiddlewareJWT = middlewares.NewAuthMiddleware(server.Config.AuthSecret)
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
