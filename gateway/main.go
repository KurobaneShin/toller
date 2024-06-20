package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the listen address of the HTTP server")
	flag.Parse()
	http.HandleFunc("/invoice", handleGetInvoice)
	logrus.Infof("gateway HTTP server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handleGetInvoice(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"invoice": "bla"})
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
