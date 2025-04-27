# Real-Time Campaign Analytics System

This project simulates a real-time ad campaign analytics platform that ingests, processes, stores, and serves campaign insights from multiple ad sources like Meta, Google, LinkedIn, and TikTok.

---

## Task Overview

Build a scalable, fault-tolerant analytics system that:
- Ingests real-time ad campaign data
- Computes metrics (CTR, ROAS, CPA)
- Stores the data in PostgreSQL
- Caches insights in Redis for low-latency APIs
- Exposes metrics via a REST API

---

## Architecture

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

## Code Flow (High-Level)

1. **main.go** initializes DB, Redis, selects ingestion mode, and starts server
2. **ingestion/** handles real API fetchers or fake data simulation
3. **processor/** computes metrics and sends them to storage
4. **storage/** manages database inserts and Redis caching
5. **api/** defines secured RESTful endpoints for fetching analytics

---

## Directory Structure

```
campaign-analytics/
├── api/                 # REST API server
├── ingestion/           # Ingestion from Meta, Google, TikTok, LinkedIn, or Fake
│   ├── meta.go
│   ├── google.go
│   ├── tiktok.go
│   ├── linkedin.go
│   ├── dispatcher.go
├── processor/           # Metric calculations
├── storage/             # PostgreSQL and Redis operations
├── models/              # Shared data models
├── Dockerfile           # App container config
├── docker-compose.yml   # Multi-service orchestration
├── init.sql             # SQL schema setup
├── go.mod / go.sum      # Go module dependencies
└── main.go              # Application entry point
```

---

## Running Instructions

### Prerequisites
- Go 1.21+
- Docker and Docker Compose

### Local Run (Simulated Fake Data)

```bash
git clone <repo-url>
cd campaign-analytics
go mod tidy
go run main.go
```

---

### Docker Compose Run

```bash
docker-compose down -v --remove-orphans
docker-compose up --build
```

Set in docker-compose:
- For fake data simulation:
  ```yaml
  - DATA_SOURCE=fake
  - ENABLED_SOURCES=
  ```
- For real data ingestion:
  ```yaml
  - DATA_SOURCE=real
  - ENABLED_SOURCES=meta,google,tiktok,linkedin
  ```

You must provide correct API credentials for real mode.

---

## API Usage

### Authentication
All API requests must include:

```bash
Authorization: Bearer <API_KEY>
```

Set `API_KEY` via environment variables.

### Endpoints

- `GET /campaign/:id/insights`

Optional query parameters:
- `from` (start date)
- `to` (end date)
- `platform` (filter by platform)

Example:

```bash
curl -H "Authorization: Bearer secret123" http://localhost:8080/campaign/cmp-42/insights?from=2024-04-01&to=2024-04-20&platform=Google
```

---

## Metrics Computed

- CTR = Clicks / Impressions
- ROAS = Revenue / Cost
- CPA = Cost / Conversions

---

## Database Schema (init.sql)

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

## Fake Data Simulation

When `DATA_SOURCE=fake`, random campaign metrics are generated every 2 seconds by `streamer.go`:

- Random impressions, clicks, cost, revenue
- Metrics processed and stored the same as real data

---

## Scaling Strategy

- Local scaling: `docker-compose up --scale app=3`
- Production scaling:
  - Deploy on Kubernetes (Horizontal Pod Autoscaler)
  - Use managed Redis (ElastiCache, Memorystore)
  - Use managed PostgreSQL (Amazon RDS, GCP Cloud SQL)
  - Add Nginx/Cloud Load Balancer

---

## Performance Benchmarking

Install [`hey`](https://github.com/rakyll/hey) and run:

```bash
hey -n 1000 -c 50 -H "Authorization: Bearer secret123" http://localhost:8080/campaign/cmp-42/insights
```

Simulates 1000 requests with 50 concurrent clients.

---

## Features Completed

| Feature                                         | Status   | Notes                                                   |
|--------------------------------------------------|----------|---------------------------------------------------------|
| Real-time ingestion from multiple platforms     | Completed | Meta, Google, TikTok, LinkedIn                          |
| Fake data simulation fallback                   | Completed | Ingestion simulation mode                              |
| Metric computation: CTR, ROAS, CPA              | Completed | In processor/aggregator.go                             |
| PostgreSQL integration                          | Completed | Inserts and queries with deduplication                 |
| Redis caching                                   | Completed | API caching using campaign IDs                         |
| REST API for insights                           | Completed | /campaign/:id/insights endpoint                        |
| API filters (date range, platform)              | Completed | Query parameters supported                             |
| Deduplication on database                       | Completed | On conflict (campaign_id, timestamp) do nothing        |
| Retry mechanism for DB inserts                  | Completed | Retry logic for transient DB errors                    |
| API authentication (Bearer token)               | Completed | Simple secure access via Authorization header         |
| Docker Compose orchestration                    | Completed | All services via docker-compose                        |
| HTTPS and Secure Deployment Notes               | Completed | Production security best practices explained          |
| Scaling strategy and performance notes          | Completed | Kubernetes, Load balancing, Horizontal scaling         |
| Modular ingestion architecture                  | Completed | Separated files for each platform ingestion            |

---

## Summary

This project demonstrates how to build a production-grade real-time ad analytics pipeline, capable of handling ingestion from real ad platforms, computing KPIs, serving REST APIs securely, caching frequently accessed data, and scaling seamlessly under increasing load.

It can serve as the foundation for marketing dashboards, real-time reporting systems, and ad optimization platforms.
