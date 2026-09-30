# movie-recommender

A small Go CLI that recommends movies with **semantic vector search**. Movie descriptions are embedded with Google's
`gemini-embedding-001` model, stored as vectors, and queried by cosine similarity against the embedding of a free-text
query such as `"bruce wayne"` or `"space exploration"`.

The same idea is implemented against two vector stores so they can be compared side by side:

| Backend | Storage | Dataset | Search |
|---|---|---|---|
| **PostgreSQL + pgvector** | `movies` table with a `vector(3072)` column | 50,000 synthetic movies (deterministically generated) | `ORDER BY description_embedding <=> $1` |
| **Redis Stack** | JSON documents under `movie:*` + `idx:movies` RediSearch index | 20 well-known movies | `FT.SEARCH` KNN query |

## How it works

![movie-recommender architecture](docs/architecture.svg)

> The diagram source is [`docs/architecture.excalidraw`](docs/architecture.excalidraw) – open it at
> <https://excalidraw.com> (File → Open) to edit, then export as SVG over `docs/architecture.svg`.

1. **Dump** – generate embeddings for every movie description and store them alongside the title and description.
2. **Recommend** – embed the query text and return the nearest movies by cosine distance
   (similarity is reported as `1 - distance`).

## Prerequisites

- Go 1.25+
- Docker (for running Postgres/pgvector and Redis Stack locally)
- A Google AI Studio API key – get one at <https://aistudio.google.com/apikey>

## Setup

```bash
git clone https://github.com/Shivam-47/movie-recommender.git
cd movie-recommender

cp .env.example .env        # then set GOOGLE_AI_API_KEY (and adjust connection details if needed)
go build -o movie-recommender .
```

The CLI loads `.env` from the current directory automatically. Variables already set in the environment take precedence.

| Variable | Used by | Description |
|---|---|---|
| `GOOGLE_AI_API_KEY` | all modes except `redis-init` | Gemini API key used for embeddings |
| `POSTGRESQL_CONNECTION_STRING` | `postgres-*` | e.g. `postgres://postgres:postgres@localhost:5432/movie-recommender` |
| `REDIS_HOST_PORT` | `redis-*` | e.g. `localhost:6379` |
| `REDIS_PASSWORD` | `redis-*` | empty for the local Docker container |

## Usage

```
movie-recommender --mode <mode> [--query "<text>"] [--top <n>]
```

| Flag | Default | Description |
|---|---|---|
| `--mode` | *(required)* | `postgres-dump`, `postgres-recommend`, `redis-init`, `redis-dump`, `redis-recommend` |
| `--query` | `bruce wayne` | Search text for the `*-recommend` modes |
| `--top` | `5` | Number of recommendations for the `*-recommend` modes |

### PostgreSQL + pgvector

```bash
docker compose -f postgres-docker-compose.yaml up -d   # creates the extension + movies table from postgres/schema.sql

./movie-recommender --mode postgres-dump                # embeds and inserts 50,000 movies
./movie-recommender --mode postgres-recommend --query "a heist inside dreams" --top 5
```

If you use your own Postgres instance instead of the compose file, run [`postgres/schema.sql`](postgres/schema.sql)
against it first.

> **Note:** `postgres-dump` makes ~500 batch embedding requests (100 descriptions each, one batch per second), so a full
> run takes a while and counts against your Gemini API quota. Change `movieCount` in
> [`postgres/dump/main.go`](postgres/dump/main.go) to load a smaller dataset. Re-running the dump inserts duplicates;
> `TRUNCATE movies;` first if you want a clean table.

### Redis Stack

```bash
docker compose -f redis-docker-compose.yaml up -d      # Redis on :6379, RedisInsight UI on http://localhost:8001

./movie-recommender --mode redis-init                   # (re)creates the idx:movies vector index
./movie-recommender --mode redis-dump                   # embeds and stores the 20 sample movies as JSON
./movie-recommender --mode redis-recommend --query "bruce wayne" --top 3
```

> `redis-init` drops any existing `idx:movies` index **together with its documents** before recreating it, so run
> `redis-dump` again afterwards.

### Example output

```
here are some recommendations for bruce wayne movies:
======
title: The Dark Knight
description: A gritty and gripping crime thriller that pits Batman against his most formidable foe, the Joker, ...
similarity score: <cosine similarity, higher is closer>
======
...
```

## Project structure

```
.
├── main.go                        # CLI entry point (cobra) – dispatches on --mode
├── postgres/
│   ├── schema.sql                 # pgvector extension + movies table
│   ├── dump/                      # generates synthetic movies, embeds in batches, inserts into Postgres
│   └── recommend/                 # embeds the query and runs a cosine-distance search
├── redis/
│   ├── index/                     # creates the RediSearch JSON + vector index (--mode redis-init)
│   ├── dump/                      # embeds 20 sample movies and stores them as JSON documents
│   └── recommend/                 # embeds the query and runs an FT.SEARCH KNN query
├── docs/
│   ├── architecture.excalidraw    # editable diagram source
│   └── architecture.svg           # rendered diagram used in this README
├── postgres-docker-compose.yaml
├── redis-docker-compose.yaml
├── .env.example
└── LICENSE
```

## Notes and limitations

- Search in both backends is **exact** (brute force): Redis uses a `FLAT` index and Postgres has no ANN index.
  pgvector's HNSW/IVFFlat indexes support at most 2,000 dimensions for `vector`, so indexing the 3,072-dim embeddings
  would require a `halfvec` expression index or a reduced output dimensionality.
- The Postgres dataset is synthetic: titles and descriptions are assembled from word lists with a fixed random seed,
  so every run produces the same 50,000 movies.

## License

[MIT](LICENSE)
