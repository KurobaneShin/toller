package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
<<<<<<< Updated upstream
	"os"
	"strconv"
=======
>>>>>>> Stashed changes

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
	fmt.Println("HTTP transport running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listenAddr, nil)
}

<<<<<<< Updated upstream
func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}

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
		if r.Method != "POST" {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}
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

=======
>>>>>>> Stashed changes
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
