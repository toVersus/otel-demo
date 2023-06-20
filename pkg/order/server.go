package order

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/toVersus/otel-demo/pkg/datastore"
	"github.com/toVersus/otel-demo/pkg/utils"
)

type Server struct {
	*http.Server

	userUrl   string
	orderAddr string

	db datastore.DB
}

func New(orderAddr, userUrl string) (*Server, error) {
	db, err := datastore.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %v", err)
	}

	return &Server{
		orderAddr: orderAddr,
		userUrl:   userUrl,
		db:        db,
	}, nil
}

func (s *Server) Setup() {
	router := http.NewServeMux()
	router.Handle("/orders", utils.LoggingMW(http.HandlerFunc(s.createOrder)))
	router.Handle("/healthz", utils.Healthz())
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	})

	srv := &http.Server{
		Addr:    s.orderAddr,
		Handler: c.Handler(router),
	}
	s.Server = srv
}

func (s *Server) Run() {
	log.Printf("Order service running at: %s", s.orderAddr)
	if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to setup http server: %v", err)
	}
}
