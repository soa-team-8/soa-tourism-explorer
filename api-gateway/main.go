package main

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

// Load configuration from config file or environment variables
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	//initConfig()

	// Create a new router
	router := mux.NewRouter()

	// Define routes and handlers
	router.HandleFunc("/encounters-service/{endpoint:.*}", handleEncountersService).Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/tours-service/{endpoint:.*}", handleToursService).Methods("GET", "POST", "PUT", "DELETE")

	// Start the server
	log.Println("API Gateway is running on port 8880")
	log.Fatal(http.ListenAndServe(":8880", router))
}

// Handler functions to proxy requests to different microservices
func handleEncountersService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	endpoint := vars["endpoint"]
	url := "http://localhost:3030" + "/" + endpoint

	proxyRequest(w, r, url)
}

func handleToursService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	endpoint := vars["endpoint"]
	url := "http://localhost:3000" + "/" + endpoint

	proxyRequest(w, r, url)
}

// Proxy function to forward requests
func proxyRequest(w http.ResponseWriter, r *http.Request, url string) {
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers and status code
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
