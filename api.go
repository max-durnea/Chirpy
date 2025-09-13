package main

import(
	"net/http"
	"fmt"
	"sync/atomic"
	"github.com/max-durnea/Server-GO/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    count := cfg.fileserverHits.Load()
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page:=fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",count)
    fmt.Fprintf(w, page)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write([]byte("Counter reset"))
}