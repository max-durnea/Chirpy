package main
import(
	_ "github.com/lib/pq"
	"net/http"
	"database/sql"
	"os"
	"fmt"
	"github.com/max-durnea/Server-GO/internal/database"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
	"time"
)

var apiCfg = apiConfig{}

//Same user struct as in internal/database package, just to have the json keys
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main(){
	//Load .env and open database connection
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("ERROR: Could not open database: %v\n",err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	apiCfg.db = dbQueries
	apiCfg.platform = os.Getenv("PLATFORM")
	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = "localhost:8080";
	fs := http.FileServer(http.Dir("."))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz",handler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /api/validate_chirp",chirp_validation)
	mux.HandleFunc("POST /api/users",apiCfg.createUserHandler)
	mux.HandleFunc("POST /admin/reset",apiCfg.resetHandler)
	server.ListenAndServe()
	
}