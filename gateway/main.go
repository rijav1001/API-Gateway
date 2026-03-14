package main

import (
	"log"
	"net/http"

	"github.com/rijav1001/API-Gateway/gateway/middleware"
	"github.com/rijav1001/API-Gateway/gateway/router"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := router.NewRouter()
	r.AddRoute("/users", []string{"http://localhost:8081"})
	r.AddRoute("/orders", []string{"http://localhost:8082"})

	rl := middleware.NewRateLimiter(5, 10) // 5 tokens/sec, max 10

	chain := middleware.Logger(logger,
		rl.Middleware(
			middleware.Auth(r),
		),
	)

	log.Println("Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", chain))
}