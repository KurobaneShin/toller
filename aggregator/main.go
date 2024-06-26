package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/KurobaneShin/tolling/types"
)

func main() {
	var (
		store    = makeStore(os.Getenv("AGG_STORE_TYPE"))
		httpAddr = os.Getenv("AGG_HTTP_ENDPOINT")
		grpcAddr = os.Getenv("AGG_GRPC_ENDPOINT")
	)
	svc := NewInvoiceAggregator(store)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGrpcTransport(grpcAddr, svc))
	}()

	log.Fatal(makeHttpTransport(httpAddr, svc))
}

func makeGrpcTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewAggregatorGrpcServer(svc))
	return server.Serve(ln)
}

func makeHttpTransport(listenAddr string, svc Aggregator) error {
	aggMetricHandler := newHTTPMetricsHandler("aggregate")
	invMetricHandler := newHTTPMetricsHandler("invoice")

	http.HandleFunc("/aggregate", aggMetricHandler.instrument(handleAggregate(svc)))
	http.HandleFunc("/invoice", invMetricHandler.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("HTTP transport running on port", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeStore(sType string) Storer {
	switch sType {
	case "memory":
		return NewMemoryStore()
	default:
		return NewMemoryStore()
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
