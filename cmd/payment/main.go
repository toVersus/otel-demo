package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/toVersus/otel-demo/pkg/payment"
	"github.com/toVersus/otel-demo/pkg/telemetry"
)

var (
	paymentAddr, userUrl string
)

func main() {
	// read the config from .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}
	paymentAddr = os.Getenv("PAYMENT_ADDR")
	userUrl = os.Getenv("USER_URL")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	tp, err := telemetry.Init("payment", "v0.1.0")
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown tracer provider: %v", err)
		}
	}()

	svc, err := payment.New(paymentAddr, userUrl)
	if err != nil {
		log.Fatalf("Failed to initialize payment service: %v", err)
	}
	svc.Setup()

	go svc.Run()

	<-ctx.Done()
	if err := svc.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shutdown payment server")
	}
}
