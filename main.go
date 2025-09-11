package main
import(
	"net/http"
	"sync/atomic"
	"fmt"   
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}



func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    count := cfg.fileserverHits.Load()
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintf(w, "Hits: %d", count)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write([]byte("Counter reset"))
}

func main(){
	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = "localhost:8080";
	var apiCfg = apiConfig{}
	fs := http.FileServer(http.Dir("."))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /healthz",handler)
	mux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
    mux.HandleFunc("POST /reset", apiCfg.resetHandler)

	server.ListenAndServe()
	
}