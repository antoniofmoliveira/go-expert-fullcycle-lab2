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
	"go.opentelemetry.io/otel/trace"

	"log"
	"net"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
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

var OtelTracer trace.Tracer
var zipkinClient *zipkinhttp.Client

func initOtel(ctx context.Context) {

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("servicoa"),
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	initOtel(ctx)

	// init zipkin
	reporter := httpreporter.NewReporter("http://zipkin-all-in-one:9411/api/v2/spans")
	localEndpoint := &model.Endpoint{ServiceName: "servicoa", IPv4: getOutboundIP(), Port: 8080}
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		log.Fatal(err)
	}
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		log.Fatal(err)
	}
	// end of initzipkin

	ctx = context.Background()
	ctx, span := OtelTracer.Start(ctx, "iniciando Servico A")
	defer span.End()

	zipkinClient, err = zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	router := http.NewServeMux()
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)
	http.Handle("/", serverMiddleware(router))
	router.HandleFunc("POST /cep", cepHandler)

	http.ListenAndServe(":8080", router)

	slog.Info("Servico A")
	select {
	case <-sigCh:
		slog.Info("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		slog.Info("Shutting down due to other reason...")
	}
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	//otel
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := OtelTracer.Start(ctx, "servicoa")
	defer span.End()

	//zipkin
	zspan := zipkin.SpanFromContext(r.Context())
	ctx = zipkin.NewContext(ctx, zspan)

	//handler
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
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header)) // IMPORTANT

	res, err := zipkinClient.Do(req)

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
}

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
