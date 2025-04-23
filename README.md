// README.md
# Real-Time Campaign Analytics System

This project simulates a real-time ad campaign analytics platform that ingests, processes, stores, and serves campaign insights from multiple ad sources like Meta, Google, LinkedIn, and TikTok.

---

## 📌 Task Overview

Build a **scalable, fault-tolerant analytics system** that:
- Ingests real-time ad campaign data
- Computes metrics (CTR, ROAS, CPA, etc.)
- Stores the data in PostgreSQL
- Caches insights in Redis for low-latency APIs
- Exposes metrics via a REST API

---

## 🏗️ Architecture

```
[Ingestion Layer] ➜ [Processor Layer] ➜ [PostgreSQL]
                                  ↓            ↑
                           [Redis Cache]    [API Layer]
```

- Language: Go
- Framework: Gin
- Data store: PostgreSQL
- Cache: Redis
- Containerization: Docker, Docker Compose

---

## 🧠 Code Flow (High-Level)

1. **main.go** initializes DB, Redis, and starts streaming + server
2. **ingestion/streamer.go** simulates incoming campaign data
3. **processor/aggregator.go** computes CTR, ROAS, CPA and stores in DB
4. **storage/db.go** handles Postgres operations
5. **storage/cache.go** handles Redis caching
6. **api/server.go** defines endpoints, adds query filtering, auth middleware

---

## 📁 Project Structure

```
campaign-analytics/
├── api/                 # Gin REST API
├── ingestion/           # Real-time simulation of ad events
├── processor/           # Metric calculation and persistence
├── storage/             # DB and Redis integration
├── models/              # Campaign struct definitions
├── Dockerfile           # Go app container setup
├── docker-compose.yml   # Multi-service definition
├── init.sql             # SQL for table creation
├── go.mod / go.sum      # Go module dependencies
└── main.go              # Application bootstrap
```

---

## 🚀 How to Run (Local)

### Prerequisites
- Go >= 1.21
- Docker & Docker Compose

### Clone and Run:
```bash
git clone <repo-url>
cd campaign-analytics
go mod tidy
go run main.go
```

---

## 🚀 How to Run (Docker Compose)

```bash
docker-compose down -v --remove-orphans
docker-compose up --build
```

Visit:
```
http://localhost:8080/campaign/cmp-42/insights
```
(*replace `cmp-42` with one seen in logs*)

---

## 🔐 API Authentication

All requests to protected endpoints must include an API key:
```http
Authorization: Bearer secret123
```

Set `API_KEY` via `docker-compose.yml` or `.env`.

Example:
```bash
curl -H "Authorization: Bearer secret123" \
  "http://localhost:8080/campaign/cmp-42/insights?from=2025-04-01&to=2025-04-20&platform=Google"
```

---

## 🔒 HTTPS & Secure Deployment

- This project runs HTTP-only (via Gin) in dev
- Use **Nginx**, **Cloudflare**, or **Kubernetes Ingress** with TLS certs for production
- Store secrets like `API_KEY` securely using vaults, `.env`, or secret managers

---

## ⚙️ Scaling Strategy

- Local scaling: `docker-compose up --scale app=3`
- Production scaling:
  - Run on Kubernetes with Horizontal Pod Autoscaling
  - Use managed Redis (ElastiCache, Memorystore)
  - Use managed Postgres (RDS, Cloud SQL)
  - Add load balancer (Nginx, GCP Load Balancer)

---

## 📊 Performance Benchmarking

Install [`hey`](https://github.com/rakyll/hey) and run:
```bash
hey -n 1000 -c 50 -H "Authorization: Bearer secret123" \
  http://localhost:8080/campaign/cmp-42/insights
```
This simulates 1000 requests with 50 concurrent clients.

---

## 📥 Fake Data Simulation

In `ingestion/streamer.go`, data is simulated like:
```go
CampaignID: fmt.Sprintf("cmp-%d", rand.Intn(100))
Platform:   random from [Meta, Google, LinkedIn, TikTok]
```
New events are streamed every 2 seconds.

---

## 🧪 Testing & Observability

### View Logs
```bash
docker-compose logs -f
```

### Inspect DB
```bash
docker exec -it postgres psql -U postgres -d campaigns
SELECT * FROM campaign_metrics LIMIT 5;
```

### Validate API
```bash
curl -H "Authorization: Bearer secret123" \
  http://localhost:8080/campaign/cmp-42/insights
```

---

## 📦 Table Schema: `init.sql`

```sql
CREATE TABLE IF NOT EXISTS campaign_metrics (
    id SERIAL PRIMARY KEY,
    campaign_id TEXT NOT NULL,
    platform TEXT NOT NULL,
    impressions INT DEFAULT 0,
    clicks INT DEFAULT 0,
    conversions INT DEFAULT 0,
    cost NUMERIC(10, 2) DEFAULT 0.00,
    revenue NUMERIC(10, 2) DEFAULT 0.00,
    timestamp TIMESTAMP NOT NULL,
    UNIQUE (campaign_id, timestamp)
);
```

---

## 📈 Metrics Computed
- **CTR** = Clicks / Impressions
- **ROAS** = Revenue / Cost
- **CPA** = Cost / Conversions

---

## ✅ Summary

This project demonstrates a production-grade campaign analytics pipeline:
- Real-time ingestion of ad metrics
- Precomputed performance KPIs (CTR, ROAS, CPA)
- Low-latency API powered by Redis
- Secure, scalable, retry-safe, deduplicated architecture

Can serve as a foundational backend for performance dashboards, real-time alerting, and marketing automation.


## ✅ Features Completed 

| Feature                                         | Status   | Notes                                                   |
|--------------------------------------------------|----------|---------------------------------------------------------|
| Real-time data ingestion                        | ✅ Done  | Simulated via Go with randomized events every 2 seconds |
| Metric computation: CTR, ROAS, CPA              | ✅ Done  | Computed in aggregator.go                              |
| PostgreSQL integration                          | ✅ Done  | Inserts, queries, and deduplication handled             |
| Redis caching                                   | ✅ Done  | Cached insights via campaign-specific keys              |
| REST API for insights                           | ✅ Done  | `/campaign/:id/insights` with filters                   |
| API filters (date range, platform)              | ✅ Done  | Supports `from`, `to`, and `platform` params            |
| Deduplication                                   | ✅ Done  | Enforced via unique constraint + insert skip            |
| Retry mechanism for DB insert                   | ✅ Done  | 3x retries with delay on DB error                       |
| API authentication (Bearer token)              | ✅ Done  | Requires `Authorization: Bearer <API_KEY>`              |
| Docker Compose support                          | ✅ Done  | Runs full stack with `docker-compose up`                |
| Table creation with `init.sql`                  | ✅ Done  | Automatically applied on Postgres init                 |
| Performance benchmarking script (`hey`)         | ✅ Done  | Docs show how to simulate 1000 requests                 |
| Scaling strategy documentation                  | ✅ Done  | Explained Docker/K8s options in README                  |
| HTTPS/secure deployment notes                   | ✅ Done  | Recommends TLS via proxy + secret handling             |
