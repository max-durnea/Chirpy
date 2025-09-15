package main

import(
	"net/http"
	"fmt"
	"sync/atomic"
	"encoding/json"
	"time"
	"strings"
	"database/sql"
	"github.com/max-durnea/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/max-durnea/Chirpy/internal/auth"
)



type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
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
    fmt.Fprintf(w,"%v",page)
}


func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request){
	type params struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	data := params{}
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	hashed_password,err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w,400,"Could not hash password")
		return
	}
	param := database.CreateUserParams{
		uuid.New(),
		time.Now(),
		time.Now(),
		data.Email,
		string(hashed_password),
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
		return
	}
	respondWithJson(w,201,user)
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

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request){
	type payload struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	var data payload
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	userID, err:=auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
		return		
	}
	//validate chirp
	if len(data.Body)>140{
		respondWithError(w,400,"Chirp is too long")
		return
	}
	joined_string := clean_string(data.Body)
	//get the user by id
	user,err := cfg.db.GetUserById(r.Context(),userID)
	if err != nil {
		respondWithError(w,400, fmt.Sprintf("%v", err))
		return
	}
	chirpParams := database.CreateChirpParams{
		uuid.New(),
		time.Now(),
		time.Now(),
		joined_string,
		user.ID,
	}
	chirpDb,err:=cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("%v", err))
		return
	}

	chirp := Chirp{
		chirpDb.ID,
		chirpDb.CreatedAt,
		chirpDb.UpdatedAt,
		chirpDb.Body,
		chirpDb.UserID,
	}
	respondWithJson(w,201,chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path
	if path == "/api/chirps"{
		dbChirps, err := cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("%v", err))
			return
		}
		chirps := make([]Chirp, len(dbChirps))
		for i, c := range dbChirps {
			chirps[i] = Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			}
		}
		respondWithJson(w, 200, chirps)
		return
	}
	prefix := "/api/chirps/"
	
	if strings.HasPrefix(path,prefix) {
		chirpIDStr := strings.TrimPrefix(path,prefix)
		id, err := uuid.Parse(chirpIDStr)

		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("%v", err))
			return
		}

		chirpDb, err := cfg.db.GetChirpById(r.Context(),id)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("%v", err))
			return
		}
		chirp := Chirp{
				chirpDb.ID,
				chirpDb.CreatedAt,
				chirpDb.UpdatedAt,
				chirpDb.Body,
				chirpDb.UserID,
		}
		respondWithJson(w,200,chirp)
	}
	
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request){
	type params struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	var data params
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	userDb, err := cfg.db.GetUserByEmail(r.Context(),data.Email)
	if err != nil {
		respondWithError(w,401,"Wrong Password/Email")
		return
	}
	err = auth.CheckPasswordHash(data.Password,userDb.HashedPassword)
	if err != nil {
		respondWithError(w,401,"Wrong Password/Email")
		return
	}
	expire := time.Duration(3600)*time.Second
	token,err:= auth.MakeJWT(userDb.ID, cfg.secret,expire)
	if err != nil {
		respondWithError(w,401,"Could not make jwt")
		return
	}
	refresh_token,err:= auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w,401,"Could not make refresh token")
		return
	}
	refresh_token_params := database.AddRefreshTokenParams{
		refresh_token,
		time.Now(),
		time.Now(),
		userDb.ID,
		time.Now().Add(60 * 24 * time.Hour),
	}
	_ , err = cfg.db.AddRefreshToken(r.Context(),refresh_token_params)
	if err != nil {
		respondWithError(w,401,"Could not add refresh token to database")
		return
	}
	resp := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
		Token:     token,
		RefreshToken: refresh_token,
	}

	respondWithJson(w, 200, resp)

}

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		respondWithError(w,400,"No Authorization Header")
		return
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix){
		respondWithError(w,400,"No Authorization Header")
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization,prefix))
	if token == ""{
		respondWithError(w,400,"No Authorization Header")
		return
	}
	tokenDb, err := cfg.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w,401,"Token does not exist")
		return
	}
	if tokenDb.RevokedAt.Valid {
		respondWithError(w, 401, "Token has been revoked")
		return
	}
	if time.Now().After(tokenDb.ExpiresAt){
		respondWithError(w,401,"Token expired")
		return
	}
	userDb,err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("%v",err))
		return
	}
	jwt_Token, err := auth.MakeJWT(userDb.ID,cfg.secret,time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w,500,"Could not generate jwt token")
		return
	}
	respondWithJson(w,200,struct{Token string `json:"token"`}{jwt_Token})


}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request){
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		respondWithError(w,400,"No Authorization Header")
		return
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix){
		respondWithError(w,400,"No Authorization Header")
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization,prefix))
	if token == ""{
		respondWithError(w,400,"No Authorization Header")
		return
	}
	tokenDb, err := cfg.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w,401,"Token does not exist")
		return
	}
	paramsToken := database.UpdateRefreshTokenParams{
		Token : tokenDb.Token,
		UpdatedAt : time.Now(),
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	_,err = cfg.db.UpdateRefreshToken(r.Context(),paramsToken)
	if err != nil {
		respondWithError(w,400,"Could not revoke token")
		return
	}
	respondWithJson(w,204,struct{}{})
}

func (cfg *apiConfig) updateHandler(w http.ResponseWriter, r *http.Request){
	// get the access token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
	}
	// get userID from token
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
	}
	userDb, err := cfg.db.GetUserById(r.Context(), userID)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
	}
	type params struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	var data params
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
	}
	hashed_pass, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
	}
	updateParams := database.UpdateUserParams{
		userID,
		data.Email,
		hashed_pass,
		time.Now(),
	}
	userDb, err = cfg.db.UpdateUser(r.Context(), updateParams)
	if err != nil {
		respondWithError(w,401,fmt.Sprintf("%v",err))
	}

	user := User{
		ID : userDb.ID,
		CreatedAt : userDb.CreatedAt,
		UpdatedAt : userDb.UpdatedAt,
		Email : userDb.Email,
	}
	respondWithJson(w,200,user)

}

