package main
import(
	"net/http"
	"sync/atomic"
	"fmt"
	"encoding/json"   
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}

func chirp_validation(w http.ResponseWriter, r *http.Request){
	type chirp struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	data := chirp{}
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w,"{\"error\":\"Something went wrong\"}")
		return
	}
	if len(data.Body)>140{
		w.WriteHeader(400)
		fmt.Fprintf(w,"{\"error\":\"Chirp is too long\"}")
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w,"{\"valid\":true}")
	return

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

func main(){
	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = "localhost:8080";
	var apiCfg = apiConfig{}
	fs := http.FileServer(http.Dir("."))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz",handler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
    mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp",chirp_validation)
	server.ListenAndServe()
	
}