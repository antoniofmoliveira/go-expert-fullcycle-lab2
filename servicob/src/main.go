package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/antoniofmoliveira/go-expert-fullcycle-lab1/src/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var OtelTracer trace.Tracer

func initOtel(ctx context.Context) {

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("servicob"),
		),
	)
	if err != nil {
		slog.Error("Resource", "failed to create resource: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	conn, err := grpc.NewClient("otel-collector:4317", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Error("gRPC", "failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		slog.Error("Trace Exporter", "failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	OtelTracer = otel.Tracer("microservice-tracer")
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	initOtel(ctx)

	ctx = context.Background()
	ctx, span := OtelTracer.Start(ctx, "iniciando Servico B")
	defer span.End()

	http.HandleFunc("/", weatherHandler)
	http.ListenAndServe(":8081", nil)

	slog.Info("Servico B")
	select {
	case <-sigCh:
		slog.Info("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		slog.Info("Shutting down due to other reason...")
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := OtelTracer.Start(ctx, "servicob")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cep := r.URL.Query().Get("cep")
	temps, status, message, err := usecase.GetWeather(ctx, cep)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, message, status)
		return
	}

	// otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header)) // !

	j, err := json.Marshal(temps)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(j))
}
