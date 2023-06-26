package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/toVersus/otel-demo/pkg/order"
	"github.com/toVersus/otel-demo/pkg/telemetry"
)

func main() {
	// read the config from .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}
	orderAddr := os.Getenv("ORDER_ADDR")
	userUrl := os.Getenv("USER_URL")
	otlpAddr := os.Getenv("OTEL_ADDR")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	closer, err := telemetry.InitTracer(ctx, otlpAddr)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer closer()

	svc, err := order.New(orderAddr, userUrl)
	if err != nil {
		log.Fatalf("Failed to initialize order service: %v", err)
	}
	svc.Setup()

	go svc.Run()

	<-ctx.Done()
	if err := svc.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shutdown order server")
	}
}
