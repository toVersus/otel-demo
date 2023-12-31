package payment

import (
	"errors"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/toVersus/otel-demo/pkg/telemetry"
	"github.com/toVersus/otel-demo/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/filters"
)

type Server struct {
	*http.Server

	paymentAddr string
	userUrl     string
}

func New(paymentAddr, userUrl string) (*Server, error) {
	return &Server{
		paymentAddr: paymentAddr,
		userUrl:     userUrl,
	}, nil
}

func (s *Server) Setup() {
	router := http.NewServeMux()
	router.Handle("/payments/transfer/id/", utils.LoggingMW(telemetry.RequestMW("payment", http.HandlerFunc(s.transferAmount))))
	router.Handle("/healthz", utils.Healthz())
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	})

	srv := &http.Server{
		Addr: s.paymentAddr,
		Handler: otelhttp.NewHandler(
			c.Handler(router), "payment",
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
	log.Printf("Payment service running at: %s", s.paymentAddr)
	if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to setup http server: %v", err)
	}
}
