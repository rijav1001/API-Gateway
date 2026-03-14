package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/orders/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"service": "order-service", "status": "ok"})
	})

	http.HandleFunc("/orders/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{
			{"id": "101", "item": "Laptop", "userId": "01"},
			{"id": "102", "item": "Phone", "userId": "02"},
		})
	})

	log.Println("Order service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil));
}