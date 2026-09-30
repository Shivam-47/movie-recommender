CREATE EXTENSION IF NOT EXISTS vector;

-- gemini-embedding-001 returns 3072-dimensional embeddings
CREATE TABLE IF NOT EXISTS movies (
    id                    BIGSERIAL PRIMARY KEY,
    title                 TEXT NOT NULL,
    description           TEXT NOT NULL,
    description_embedding vector(3072) NOT NULL
);
