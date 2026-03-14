package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/users/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"service": "user-service", "status": "ok"})
	})

	http.HandleFunc("/users/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"id": "01", "name": "John"},
			{"id": "02", "name": "Anderson"},
		})
	})

	log.Println("User service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil));
}