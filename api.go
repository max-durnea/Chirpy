package main

import(
	"net/http"
	"fmt"
	"sync/atomic"
	"encoding/json"
	"time"
	"github.com/max-durnea/Server-GO/internal/database"
	"github.com/google/uuid"
	
)



type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
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


func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request){
	type params struct{
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	data := params{}
	err := decoder.Decode(&data)
	if err != nil {

		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	param := database.CreateUserParams{
		uuid.New(),
		time.Now(),
		time.Now(),
		data.Email,
	}
	dbUser,err := cfg.db.CreateUser(r.Context(),param)
	user := User{
        ID:        dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email:     dbUser.Email,
    }

	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
	}
	repsondWithJson(w,201,user)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request){
	if cfg.platform != "dev"{
		respondWithError(w,403,"")
		return
	}
	err:=cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	w.WriteHeader(200)
	
}