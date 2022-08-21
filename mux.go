package main

import (
	"fmt"
	"net/http"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		rMsg := fmt.Sprintf(`{"status": "ok", "path": "%s"}`, r.URL.Path[1:])
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(rMsg))
	})
	return mux
}
