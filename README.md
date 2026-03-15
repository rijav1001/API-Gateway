# API Gateway

A production-grade API Gateway built in Go, deployed on Kubernetes. Built as a portfolio project to demonstrate backend engineering and DevOps skills.

## Architecture
```
Client (curl / Postman)
        │
        ▼
┌─────────────────────┐
│     API Gateway      │  :8080
│                     │
│  • JWT Auth         │
│  • Rate Limiting    │
│  • Round-Robin LB   │
│  • Structured Logs  │
│  • Live Dashboard   │
└────────┬────────────┘
         │
    ┌────┴────┐
    ▼         ▼
User-Service  Order-Service
  :8081         :8082
```

## Features

- **Routing** — Path-based routing to backend microservices
- **JWT Authentication** — Validates Bearer tokens on all protected routes
- **Rate Limiting** — Token bucket algorithm, per-IP (5 req/sec, max 10)
- **Load Balancing** — Round-robin across multiple service replicas
- **Structured Logging** — JSON logs via Uber Zap (method, path, status, latency)
- **Live Dashboard** — Real-time request stats, route hits, latency, rate limit tracking
- **Kubernetes Ready** — ConfigMap, Secrets, Deployments, Services on Minikube
- **Docker Compose** — Single command local development setup

## Tech Stack

- **Go** — Gateway and backend services
- **Docker** — Containerization
- **Kubernetes (Minikube)** — Container orchestration
- **Uber Zap** — Structured logging
- **golang-jwt** — JWT authentication

## Project Structure
```
go-api-gateway/
├── Dockerfile.gateway
├── docker-compose.yml
├── gateway/
│   ├── main.go
│   ├── router/          # Path-based routing
│   ├── middleware/      # JWT auth, rate limiting, logging
│   ├── loadbalancer/    # Round-robin load balancer
│   └── dashboard/       # Live stats dashboard
├── services/
│   ├── user-service/    # Mock user service (:8081)
│   └── order-service/   # Mock order service (:8082)
└── k8s/
    ├── configmap.yaml
    ├── gateway.yaml
    ├── user-service.yaml
    └── order-service.yaml
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker Desktop
- Minikube
- kubectl

### Run with Docker Compose (Local)
```bash
docker compose up --build
```

Gateway available at `http://localhost:8080`

### Run on Kubernetes (Minikube)
```bash
# Start Minikube
minikube start

# Point Docker to Minikube's daemon (PowerShell)
minikube docker-env --shell powershell | Invoke-Expression

# Build images
docker build -f Dockerfile.gateway -t api-gateway:v1 .
docker build -t user-service:v1 services/user-service/
docker build -t order-service:v1 services/order-service/

# Deploy
kubectl apply -f k8s/

# Get Minikube IP
minikube ip
```

Gateway available at `http://<MINIKUBE-IP>:30080`

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /users/list | ✅ JWT | List all users |
| GET | /users/health | ✅ JWT | User service health |
| GET | /orders/list | ✅ JWT | List all orders |
| GET | /orders/health | ✅ JWT | Order service health |
| GET | /dashboard | ❌ | Live stats dashboard |
| GET | /dashboard/stats | ❌ | Raw stats JSON |

## Generating a JWT Token

1. Go to [jwt.io](https://jwt.io)
2. Set payload: `{"sub": "testuser"}`
3. Set secret: `super-secret-key`
4. Copy the generated token

## Testing
```bash
# Unauthorized request
curl http://localhost:8080/users/list

# Authorized request
curl -H "Authorization: Bearer <token>" http://localhost:8080/users/list

# Trigger rate limiting
for i in {1..20}; do curl -H "Authorization: Bearer <token>" http://localhost:8080/users/list; done
```

## Dashboard

Open `http://localhost:8080/dashboard` to see live metrics:

- Total requests served
- Average latency (ms)
- Rate limit hits
- Per-route request counts