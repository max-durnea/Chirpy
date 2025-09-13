package main

import(
	"net/http"
	"sync/atomic"
	"fmt"
	"encoding/json"
	"strings"
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
		respondWithError(w,400,"Something went wrong")
		return
	}
	if len(data.Body)>140{
		respondWithError(w,400,"Chirp is too long")
		return
	}
	joined_string := clean_string(data.Body)
	repsondWithJson(w,200,struct{CleanedBody string `json:"cleaned_body"`}{joined_string})

}
func clean_string(msg string) (cleaned_string string){
	words := strings.Split(msg," ")
	banned_words := map[string]struct{}{"kerfuffle":{},"sharbert":{},"fornax":{}}
	const censor  = "****"
	for i,word := range words {
		if _, ok := banned_words[strings.ToLower(word)]; ok {
			words[i] = censor
		}
	}
	cleaned_string = strings.Join(words," ")
	return cleaned_string
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