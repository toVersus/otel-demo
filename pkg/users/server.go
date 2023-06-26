package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/toVersus/otel-demo/pkg/datastore"
	"github.com/toVersus/otel-demo/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/filters"
)

type Server struct {
	*http.Server

	userAddr string

	db datastore.DB
}

func New(userAddr string) (*Server, error) {
	db, err := datastore.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %v", err)
	}

	return &Server{
		userAddr: userAddr,
		db:       db,
	}, nil
}

func (s *Server) Setup() {
	router := http.NewServeMux()
	router.Handle("/users", http.HandlerFunc(s.createUser))
	router.Handle("/users/", http.HandlerFunc(s.manageUser))
	router.Handle("/healthz", utils.Healthz())
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	})

	srv := &http.Server{
		Addr: s.userAddr,
		Handler: otelhttp.NewHandler(
			c.Handler(router), "users",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			// Ignore healthz endpoint from tracing
			otelhttp.WithFilter(filters.All(
				filters.Not(filters.Path("/healthz")),
			)),
		),
	}
	s.Server = srv
}

func (s *Server) Run() {
	log.Printf("Users service running at: %s", s.userAddr)
	if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to setup http server: %v", err)
	}
}
