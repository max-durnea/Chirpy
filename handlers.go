package main

import(
	"net/http"
	"encoding/json"
)

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
	respondWithJson(w,200,struct{CleanedBody string `json:"cleaned_body"`}{joined_string})

}

