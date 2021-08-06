package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var acceptTypes = []string{"text/plain", "application/json"}

// GET rejects HTTP methods apart from GET and HEAD
func GET(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead && r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func contentType(r *http.Request) string {
	contentType := r.Header["Content-Type"]
	if len(contentType) == 0 {
		return "text/plain"
	}
	parts := strings.Split(contentType[0], ";")
	return parts[0]
}

func respond(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	switch contentType(r) {
	case "text/plain":
		w.WriteHeader(code)
		fmt.Fprintln(w, v)
	case "application/json":
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(v); err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func main() {
	http.HandleFunc("/", GET(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "endpoints: /version /duration")
	}))
	http.HandleFunc("/version", GET(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "1.0")
	}))
	http.HandleFunc("/duration", GET(func(w http.ResponseWriter, r *http.Request) {
		bt, err := readBootTime(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}
		respond(w, r, http.StatusOK, bt)
	}))

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	fmt.Fprintf(os.Stderr, "Server ready at %s, endpoints: /version and /duration\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
