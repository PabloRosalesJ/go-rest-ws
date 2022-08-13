package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/PabloRosalesJ/go/rest-ws/handlers"
	"github.com/PabloRosalesJ/go/rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error to load .env:")
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")
	SECRET := os.Getenv("TWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal("Cant create server")
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
}
