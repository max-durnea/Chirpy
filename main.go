package main
import(
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}



func main(){
	mux := http.NewServeMux()
	server := &http.Server{}
	server.Handler = mux
	server.Addr = "localhost:8080";
	mux.Handle("/app/",http.StripPrefix("/app/",http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz",handler)
	server.ListenAndServe()
	
}