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



# Campaign Analytics Conversational Bot

This project extends the real-time Campaign Analytics system with a lightweight conversational agent.  
Users can query campaign insights using free-form prompts like:

- "What’s the ROAS for my latest Google campaign?"
- "Show me CTR trends for last week’s Meta campaigns."
- "How much did I spend on LinkedIn ads this quarter?"

---

## Task Overview

- Accept and parse free-form user prompts
- Map them into structured queries (intent and filters)
- Fetch campaign analytics from backend APIs
- Format responses in conversational style
- Support vector-based matching for fuzzy queries
- Ready for future integration with Slack, WhatsApp, etc.

---

## Architecture

[User Prompt] ➜ [Bot Server (Gin)] ↓ [Prompt Parser (Intent/Filters)] ↓ [Vector Search (Optional for fuzzy matching)] ↓ [Analytics API (Campaign App)] ↓ [Format into Conversational Response]


- Bot Language: Go
- Framework: Gin
- Backend: Campaign Analytics App
- Optional: Embedding generation and vector search (pgvector)

---

## Code Flow (High-Level)

1. **bot_main.go**:  
   Starts the Bot server on port `8081`.
2. **handler.go**:  
   Receives the `/prompt` POST request.
3. **parser.go**:  
   Parses the natural language prompt into structured **intent** and **filters**.
4. **vector_search.go** *(optional)*:  
   Uses vector embedding similarity to better understand complex prompts.
5. **client.go**:  
   Queries the Campaign Analytics API to fetch insights.
6. **formatter.go**:  
   Formats the API result into human-readable response text.
7. Returns a final JSON response.

---

## Directory Structure

campaign-analytics/ ├── bot/ │ ├── bot_main.go # Bot server entry point │ ├── handler.go # Receives and processes user prompt │ ├── parser.go # Parses prompts into structured intents │ ├── client.go # Fetches data from Campaign Analytics API │ ├── formatter.go # Formats the API response into text │ ├── vector_search.go # (Optional) Vector similarity matching │ ├── embedding.go # (Optional) Prompt embedding generator ├── Dockerfile.bot # Separate Dockerfile to build bot ├── docker-compose.yml # Multi-service orchestration


---

## Running Instructions

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Campaign Analytics App (`app`) running on port `8080`
- PostgreSQL with `pgvector` extension installed (only if vector search is used)

---

### Docker Compose Run

```bash
docker-compose down -v --remove-orphans
docker-compose up --build
```

Bot will be available at:

http://localhost:8081/prompt

API Usage (Bot Server)
Endpoint
POST /prompt

Request Body:


Copy
```
{
  "prompt": "Show me CTR trends for last week’s Meta campaigns."
}
```


Example CURL:

```bash

curl -X POST http://localhost:8081/prompt \
-H "Content-Type: application/json" \
-d '{"prompt": "What is the ROAS for my latest Google campaign?"}'
```

Example Prompts and Responses

User Prompt	Bot Response Example
"What is the ROAS for my latest Google campaign?"	"The ROAS for your campaign is 7.82."
"Show CTR for last week's Meta ads"	"The CTR for your campaign is 12.45%."
"Spend for LinkedIn this quarter"	"You spent $1450.00 on LinkedIn this quarter."
Vector Search and Embedding Support
vector_search.go:
Uses cosine similarity over text embeddings to better match fuzzy or unstructured prompts.

embedding.go:
Encodes prompts into vector embeddings using simple TF-IDF or external APIs (if added).

Important Note:
Currently a basic implementation; can later upgrade to OpenAI Embeddings, HuggingFace transformers, or pgvector extension in Postgres for production use.

Future Extensibility
Integrate with Slack Bot, WhatsApp Bot

Upgrade parsing from keyword-matching → proper NLP (spaCy, transformers)

Use semantic search (vector DB) instead of simple matching

Improve prompt-to-SQL mapping automatically

Add authentication for incoming user prompts

Features Completed

Feature	Status	Notes
Free-form prompt ingestion	Completed	Supports basic campaign-related questions
Intent and filter extraction	Completed	Metric + Platform + Dates parsing
Backend Analytics API integration	Completed	Secure call with Bearer Token
Conversational response formatting	Completed	User-friendly text responses
Dockerized bot server	Completed	Runs via separate container
Vector embedding and matching (optional)	Completed	Basic matching via cosine similarity
Final Summary
The Bot server adds a natural conversational layer on top of the real-time Campaign Analytics system.
It turns structured campaign metrics into simple, intuitive human dialogue — and is designed to easily integrate with any chat-based interface in the future.

