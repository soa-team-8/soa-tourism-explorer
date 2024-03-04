package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func loadRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Use(loggerMiddleware)

	router.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return router
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		fmt.Printf(
			"[%s] %s %s %v\n",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
