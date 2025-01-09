package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"users-service/api"
	"users-service/config"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitializeTracer(endpoint string) (*trace.TracerProvider, error) {
	// Configure the Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create a TracerProvider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("WedoLow"),
		)),
	)

	// Register the global TracerProvider
	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	cfg := config.NewConfig()
	server := api.NewServer(cfg)

	// Initialize the tracer
	tp, err := InitializeTracer(cfg.JaegerEndpoint)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
