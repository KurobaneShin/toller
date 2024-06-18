package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"github.com/KurobaneShin/tolling/types"
)

func main() {
	httpAddr := flag.String("httpAddr", ":3000", "the listen address of the HTTP server")
	grpcAddr := flag.String("grpcAddr", ":3001", "the listen address of the HTTP server")
	flag.Parse()
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGrpcTransport(*grpcAddr, svc))
	}()

	log.Fatal(makeHttpTransport(*httpAddr, svc))
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
	fmt.Println("HTTP transport running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	return http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuId := r.URL.Query().Get("obu")
		if len(obuId) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu query param"})
			return
		}

		obuID, err := strconv.Atoi(obuId)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "obuId is not a number"})
			return
		}

		inv, err := svc.CalculateInvoice(obuID)
		if err != nil {

			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, inv)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
