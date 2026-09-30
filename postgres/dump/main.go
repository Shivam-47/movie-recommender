package dump

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/jackc/pgx/v5/pgxpool"
	pgv "github.com/pgvector/pgvector-go"
	"google.golang.org/api/option"
)

var pool *pgxpool.Pool

var em *genai.EmbeddingModel

const movieCount = 50000
const batchSize = 100

type embeddingResult struct {
	titles     []string
	descs      []string
	embeddings [][]float32
}

func Main() {

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

	movieData := generateMovies(movieCount)

	titles := make([]string, 0, len(movieData))
	descriptions := make([]string, 0, len(movieData))
	for t, d := range movieData {
		titles = append(titles, t)
		descriptions = append(descriptions, d)
	}

	total := len(titles)
	resultCh := make(chan embeddingResult)

	// producer: fetch embeddings in batches, send to channel
	go func() {
		defer close(resultCh)
		for i := 0; i < total; i += batchSize {
			end := min(i+batchSize, total)

			batchTitles := titles[i:end]
			batchDescs := descriptions[i:end]

			embeddings, err := getBatchEmbeddings(batchDescs)
			if err != nil {
				log.Fatal(err)
			}

			resultCh <- embeddingResult{
				titles:     batchTitles,
				descs:      batchDescs,
				embeddings: embeddings,
			}

			time.Sleep(time.Second * 1)
		}
	}()

	// consumer: write to DB as soon as each batch arrives
	saved := 0
	for result := range resultCh {
		for j, embedding := range result.embeddings {
			err := loadData(result.titles[j], result.descs[j], embedding)
			if err != nil {
				log.Fatal(err)
			}
		}
		saved += len(result.embeddings)
		fmt.Printf("processed %d/%d\n", saved, total)
	}

	fmt.Println("data loading complete")
}

func loadData(title, description string, descriptionEmbedding []float32) error {

	_, err := pool.Exec(context.Background(), "INSERT INTO movies (title, description, description_embedding) VALUES ($1, $2, $3)", title, description, pgv.NewVector(descriptionEmbedding))

	return err
}

func getBatchEmbeddings(inputs []string) ([][]float32, error) {

	ctx := context.Background()

	batch := em.NewBatch()
	for _, input := range inputs {
		batch.AddContent(genai.Text(input))
	}

	res, err := em.BatchEmbedContents(ctx, batch)
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(res.Embeddings))
	for i, e := range res.Embeddings {
		embeddings[i] = e.Values
	}

	return embeddings, nil
}
