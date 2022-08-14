package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/PabloRosalesJ/go/rest-ws/database"
	"github.com/PabloRosalesJ/go/rest-ws/repository"
	"github.com/gorilla/mux"
)

/* Will be used to handle multiple server instances */
type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	log.Println("Starting server on port", b.Config().Port)

	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	repository.SetRepository(repo)

	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("prot is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("secrete is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("databaseURL is required")
	}

	return &Broker{
		config: config,
		router: mux.NewRouter(),
	}, nil
}
