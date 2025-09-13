package main
import(
	_ "github.com/lib/pq"
	"net/http"
	"database/sql"
	"os"
	"fmt"
	"github.com/max-durnea/Server-GO/internal/database"
	"github.com/joho/godotenv"
)

var apiCfg = apiConfig{}

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

	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = "localhost:8080";
	fs := http.FileServer(http.Dir("."))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz",handler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
    mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp",chirp_validation)
	server.ListenAndServe()
	
}