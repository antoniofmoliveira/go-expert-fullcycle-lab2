package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type ceprequest struct {
	Cep string `json:"cep"`
}

func (c ceprequest) validate() error {
	r, _ := regexp.Compile("^[0-9]{8}$")
	if !r.MatchString(c.Cep) {
		return errors.New("invalid zipcode")
	}

	return nil
}

func main() {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("servicoa"),
		),
	)
	if err != nil {
		slog.Error("Resource", "failed to create resource: %w", err)
	}
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// conn, err := grpc.DialContext(ctx, "otel-collector:4317",
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	// grpc.WithBlock(),
	// )

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

	tracer := otel.Tracer("microservice-tracer")

	ctx = context.Background()
	ctx, span := tracer.Start(ctx, "iniciando Servico A")
	defer span.End()

	http.HandleFunc("POST /cep", func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.HeaderCarrier(r.Header)
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		// tracer := otel.Tracer("microservice-tracer")

		ctx, span := tracer.Start(ctx, "servicoa")
		defer span.End()

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header)) // !

		body, error := io.ReadAll(r.Body)
		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var cep ceprequest
		if error := json.Unmarshal(body, &cep); error != nil {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if error := cep.validate(); error != nil {
			http.Error(w, error.Error(), http.StatusUnprocessableEntity)
			return
		}
		url := "http://servicob:8081/?cep={{cep}}"
		url = strings.Replace(url, "{{cep}}", cep.Cep, 1)
		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		slog.Info("status", "code", res.StatusCode)
		switch res.StatusCode {
		case http.StatusOK:
			body, err := io.ReadAll(res.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer res.Body.Close()
			sbody := string(body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(sbody))
		case http.StatusNotFound:
			http.Error(w, "not found", http.StatusNotFound)
		case http.StatusUnprocessableEntity:
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})

	// go func() {
	http.ListenAndServe(":8080", nil)
	// }()

	slog.Info("Servico A")
	select {
	case <-sigCh:
		slog.Info("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		slog.Info("Shutting down due to other reason...")
	}
}

// func cepHandler(w http.ResponseWriter, r *http.Request) {
// 	carrier := propagation.HeaderCarrier(r.Header)
// 	ctx := r.Context()
// 	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

// 	tracer := otel.Tracer("microservice-tracer")

// 	ctx, span := tracer.Start(ctx, "servicoa")
// 	defer span.End()

// 	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header)) // !

// 	body, error := io.ReadAll(r.Body)
// 	if error != nil {
// 		http.Error(w, error.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var cep ceprequest
// 	if error := json.Unmarshal(body, &cep); error != nil {
// 		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
// 		return
// 	}
// 	if error := cep.validate(); error != nil {
// 		http.Error(w, error.Error(), http.StatusUnprocessableEntity)
// 		return
// 	}
// 	url := "http://servicob:8081/?cep={{cep}}"
// 	url = strings.Replace(url, "{{cep}}", cep.Cep, 1)
// 	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	slog.Info("status", "code", res.StatusCode)
// 	switch res.StatusCode {
// 	case http.StatusOK:
// 		body, err := io.ReadAll(res.Body)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		defer res.Body.Close()
// 		sbody := string(body)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(sbody))
// 	case http.StatusNotFound:
// 		http.Error(w, "not found", http.StatusNotFound)
// 	case http.StatusUnprocessableEntity:
// 		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
// 	default:
// 		http.Error(w, "internal server error", http.StatusInternalServerError)
// 	}
// }
