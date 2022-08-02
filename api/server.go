package api

import (
	"database/sql"
	"encoding/json"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type server struct {
	db     *sql.DB
	router *mux.Router
	debug  bool
}

func NewServer(db *sql.DB) *server {
	s := server{db: db, router: mux.NewRouter()}
	return &s
}

func (s *server) cors() {

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowedMethods([]string{"POST", "DELETE", "GET", "PUT", "OPTIONS"}),
	)
	s.router.Use(cors)
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func (s *server) respond(w http.ResponseWriter, v interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Println(err)
	}
}

func (s *server) LogDebug(v ...interface{}) {
	if s.debug {
		log.Println(v...)
	}
}

func (s *server) getParams(r *http.Request) map[string]string {
	params := mux.Vars(r)

	return params
}

func (s *server) getQueryParams(r *http.Request) url.Values {
	return r.URL.Query()
}

func (s *server) error(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)
	s.respond(w, struct {
		Error string `json:"error"`
	}{err.Error()}, statusCode)
}

func (s *server) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *server) registerRoute(path string, handler http.HandlerFunc, methods []string, middleware ...func(handlerFunc http.HandlerFunc) http.HandlerFunc) {
	for _, mw := range middleware {
		handler = mw(handler)
	}

	handler = s.middlewareJson(handler)

	if s.debug {
		handler = s.middlewareAccessLog(handler)
	}

	s.router.HandleFunc(path, handler).Methods(methods...)
}

func (s *server) HTTPServe(port string) {
	s.debug, _ = strconv.ParseBool(os.Getenv("API_DEBUG"))
	s.routes()
	s.cors()
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, s.router))
}

func (s *server) LambdaServe(mux *gorillamux.GorillaMuxAdapter) {
	s.routes()
	s.cors()
	mux = gorillamux.New(s.router)
}
