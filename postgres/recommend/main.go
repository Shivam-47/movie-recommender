package recommend

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/jackc/pgx/v5/pgxpool"
	pgv "github.com/pgvector/pgvector-go"
	"google.golang.org/api/option"
)

var pool *pgxpool.Pool
var em *genai.EmbeddingModel

const searchQuery = "SELECT title, description, 1 - (description_embedding <=> $1) as similarity FROM movies ORDER BY description_embedding <=> $1 LIMIT $2"

func Main(queryString string, topN int) {

	connString := os.Getenv("POSTGRESQL_CONNECTION_STRING")
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")

	ctx := context.Background()
	var err error

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	em = client.EmbeddingModel("gemini-embedding-001")

	pool, err = pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queryEmbedding, err := getEmbedding(queryString)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := pool.Query(ctx, searchQuery, pgv.NewVector(queryEmbedding), topN)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var title string
	var description string
	var score float32

	fmt.Println("here are some recommendations for", queryString, "movies:\n======")

	for rows.Next() {
		err = rows.Scan(&title, &description, &score)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("title:", title)
		fmt.Println("description:", description)
		fmt.Println("similarity score:", score)

		fmt.Println("======")
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func getEmbedding(input string) ([]float32, error) {

	ctx := context.Background()

	res, err := em.EmbedContent(ctx, genai.Text(input))

	if err != nil {
		return nil, err
	}

	return res.Embedding.Values, nil

}
