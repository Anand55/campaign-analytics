// README.md
# Real-Time Campaign Analytics System

This project simulates a real-time ad campaign analytics platform that ingests, processes, stores, and serves campaign insights from multiple ad sources like Meta, Google, LinkedIn, and TikTok.

---

## 📌 Task Overview

Build a **scalable, fault-tolerant analytics system** that:
- Ingests real-time ad campaign data
- Computes metrics (CTR, ROAS, etc.)
- Stores the data in PostgreSQL
- Caches insights in Redis for low-latency APIs
- Exposes metrics via a REST API

---

## 🏗️ Architecture

```
[Ingestion Layer]
   ↳ Simulates real-time campaign data (every 2 seconds)
       ↓
[Processor Layer]
   ↳ Calculates CTR, ROAS, logs, and stores in PostgreSQL
       ↓
[Storage Layer]
   ↳ PostgreSQL for persistence
   ↳ Redis for fast lookup / caching
       ↓
[API Layer]
   ↳ Exposes insights via GET /campaign/:id/insights
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
3. **processor/aggregator.go** computes CTR, ROAS, and saves to DB
4. **storage/db.go** handles Postgres connection and inserts
5. **storage/cache.go** integrates Redis for caching
6. **api/server.go** defines REST endpoints and uses DB + Redis

---

## 📁 Project Structure

```
campaign-analytics/
├── api/                 # Gin REST API
├── ingestion/           # Real-time simulation of ad events
├── processor/           # CTR & ROAS calculator + DB persistence
├── storage/             # DB (Postgres) and Cache (Redis) logic
├── models/              # Structs and campaign schema
├── Dockerfile           # Go service container setup
├── docker-compose.yml   # Multi-container orchestration
├── init.sql             # Bootstrap SQL to create campaign_metrics table
├── go.mod / go.sum      # Module dependencies
└── main.go              # Entrypoint, orchestrates everything
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

### 1. Build and Start All Services
```bash
docker-compose down -v --remove-orphans
docker-compose up --build
```

### 2. Access REST API
Visit:
```
http://localhost:8080/campaign/cmp-42/insights
```
(*use actual `cmp-XX` from the logs*)

---

## 📥 Fake Data Simulation

- Located in `ingestion/streamer.go`
- Uses `rand` to generate fake campaigns every 2 seconds:
  ```go
  CampaignID: fmt.Sprintf("cmp-%d", rand.Intn(100)),
  Platform:   platforms[rand.Intn(len(platforms))],
  Impressions, Clicks, Cost, Revenue => randomized
  ```

---

## 🧪 Testing & Observability

### View Logs
```bash
docker-compose logs -f
```

### Query DB
```bash
docker exec -it postgres psql -U postgres -d campaigns
SELECT * FROM campaign_metrics LIMIT 5;
```

### Validate API
```bash
curl http://localhost:8080/campaign/cmp-42/insights
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
    timestamp TIMESTAMP NOT NULL
);
```

---

## 📈 Metrics Computed
- **CTR** = Clicks / Impressions
- **ROAS** = Revenue / Cost

---

## ✅ Summary
This app mimics how an ad analytics system would work in production — streaming data, persisting events, caching hot queries, and exposing insights — all within a Dockerized microservice setup.

Use it to demonstrate:
- Real-time ingestion
- REST APIs with caching
- Metric computation
- Distributed architecture fundamentals
