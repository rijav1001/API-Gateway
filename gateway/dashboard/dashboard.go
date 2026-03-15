package dashboard

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Stats struct {
	mu            	sync.Mutex
	TotalRequests 	int64            `json:"total_requests"`
	RouteHits     	map[string]int64 `json:"route_hits"`
	RateLimitHits 	int64            `json:"rate_limit_hits"`
	TotalLatencyMs 	float64          `json:"avg_latency_ms"`
	requestCount  	int64
}

var Global = &Stats{
	RouteHits: make(map[string]int64),
}

func (s *Stats) RecordRequest(path string, latencyMs float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalRequests++
	s.requestCount++
	s.RouteHits[path]++
	s.TotalLatencyMs = (s.TotalLatencyMs * float64(s.requestCount-1) + latencyMs) / float64(s.requestCount)
}

func (s *Stats) RecordRateLimit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RateLimitHits++
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/dashboard/stats" {
		w.Header().Set("Content-Type", "application/json")
		Global.mu.Lock()
		json.NewEncoder(w).Encode(Global)
		Global.mu.Unlock()
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashboardHTML))
}

var dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>API Gateway Dashboard</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: 'Segoe UI', sans-serif; background: #0f172a; color: #e2e8f0; min-height: 100vh; padding: 2rem; }
    h1 { font-size: 1.8rem; margin-bottom: 2rem; color: #38bdf8; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
    .card { background: #1e293b; border-radius: 12px; padding: 1.5rem; border: 1px solid #334155; }
    .card h2 { font-size: 0.85rem; color: #94a3b8; margin-bottom: 0.5rem; text-transform: uppercase; letter-spacing: 0.05em; }
    .card .value { font-size: 2.2rem; font-weight: 700; color: #38bdf8; }
    .routes { background: #1e293b; border-radius: 12px; padding: 1.5rem; border: 1px solid #334155; }
    .routes h2 { font-size: 0.85rem; color: #94a3b8; margin-bottom: 1rem; text-transform: uppercase; letter-spacing: 0.05em; }
    .route-row { display: flex; justify-content: space-between; padding: 0.6rem 0; border-bottom: 1px solid #334155; }
    .route-row:last-child { border-bottom: none; }
    .badge { background: #0ea5e9; color: white; padding: 0.2rem 0.7rem; border-radius: 999px; font-size: 0.85rem; font-weight: 600; }
  </style>
</head>
<body>
  <h2>⚡ API Gateway Dashboard</h2>
  <div class="grid">
    <div class="card"><h3>Total Requests</h3><div class="value" id="total">-</div></div>
    <div class="card"><h3>Avg Latency</h3><div class="value" id="latency">-</div></div>
    <div class="card"><h3>Rate Limit Hits</h3><div class="value" id="ratelimit">-</div></div>
  </div>
  <div class="routes">
    <h3>Route Hits</h3>
    <div id="routes">Loading...</div>
  </div>
  <script>
    async function refresh() {
      const res = await fetch('/dashboard/stats');
      const data = await res.json();
      document.getElementById('total').textContent = data.total_requests;
      document.getElementById('latency').textContent = data.avg_latency_ms.toFixed(2) + ' ms';
      document.getElementById('ratelimit').textContent = data.rate_limit_hits;
      const routesDiv = document.getElementById('routes');
      routesDiv.innerHTML = '';
      for (const [route, hits] of Object.entries(data.route_hits || {})) {
        routesDiv.innerHTML += '<div class="route-row"><span>' + route + '</span><span class="badge">' + hits + '</span></div>';
      }
    }
    refresh();
    setInterval(refresh, 2000);
  </script>
</body>
</html>`

// Middleware to record stats
func StatsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dashboard" || r.URL.Path == "/dashboard/stats" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		Global.RecordRequest(r.URL.Path, float64(time.Since(start).Milliseconds()))
	})
}