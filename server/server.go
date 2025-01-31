package server

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

type Server struct {
	*http.ServeMux
}

type Route struct {
	Path    string
	Handler http.Handler
}

// newCorsHandler creates a new CORS handler
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Connect-Protocol-Version,Accept,Authorization,Content-Type,X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ConnectServer connects a server to a port with multiple routes
func (s *Server) ConnectServer(port string, routes []Route) error {
	r := mux.NewRouter()
	r.Use(cors)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Add each route to the router
	for _, route := range routes {
		r.NewRoute().PathPrefix(route.Path).Handler(route.Handler)
	}

	// Create a server object with a custom error log
	server := &http.Server{
		Addr:     port,
		Handler:  r,
		ErrorLog: log.New(os.Stderr, "server error: ", log.Lshortfile),
	}

	// Start the server
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed")
		} else {
			log.Printf("Error starting server: %s\n", err)
		}
		return err
	}

	return nil
}
