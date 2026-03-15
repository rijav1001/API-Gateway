package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rijav1001/API-Gateway/gateway/dashboard"
	"github.com/rijav1001/API-Gateway/gateway/middleware"
	"github.com/rijav1001/API-Gateway/gateway/router"

	"go.uber.org/zap"
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := router.NewRouter()

	userService := getEnv("USER_SERVICE_URL", "http://localhost:8081")
	orderService := getEnv("ORDER_SERVICE_URL", "http://localhost:8082")
	if err := r.AddRoute("/users", []string{userService}); err != nil {
		log.Fatal("Failed to add users route:", err)
	}
	if err := r.AddRoute("/orders", []string{orderService}); err != nil {
		log.Fatal("Failed to add orders route:", err)
	}

	rl := middleware.NewRateLimiter(5, 10) // 5 tokens/sec, max 10

	mux := http.NewServeMux()
	mux.HandleFunc("/dashboard", dashboard.Handler)
	mux.HandleFunc("/dashboard/stats", dashboard.Handler)
	mux.Handle("/", r)

	chain := dashboard.StatsMiddleware(
		middleware.Logger(logger,
			rl.Middleware(
				middleware.Auth(r),
			),
		),
	)

	// Dashboard doesn't need auth
	finalMux := http.NewServeMux()
	finalMux.HandleFunc("/dashboard", dashboard.Handler)
	finalMux.HandleFunc("/dashboard/stats", dashboard.Handler)
	finalMux.Handle("/", chain)

	log.Println("Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", finalMux))
}