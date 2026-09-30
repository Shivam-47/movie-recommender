package recommend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/generative-ai-go/genai"
	"github.com/redis/rueidis"
	"google.golang.org/api/option"
)

const indexName = "idx:movies"
const vectorScoreAttribute = "vector_score"
const defaultRedisUsername = "default"

var em *genai.EmbeddingModel
var client rueidis.Client

const queryTemplate = "(*)=>[KNN %d @vector $v AS %s]"

func Main(queryString string, topN int) {

	redisHostPort := os.Getenv("REDIS_HOST_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")

	var err error

	client, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress:       []string{redisHostPort},
		Username:          defaultRedisUsername,
		Password:          password,
		ForceSingleClient: true,
		DisableCache:      true,
	})

	if err != nil {
		log.Fatal("conn error - ", err)
	}
	defer client.Close()

	genaiC, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	defer genaiC.Close()

	em = genaiC.EmbeddingModel("gemini-embedding-001")

	embeddedQ, err := getEmbedding(queryString)
	if err != nil {
		log.Fatal(err)
	}

	blob := rueidis.VectorString32(embeddedQ)

	searchQuery := fmt.Sprintf(queryTemplate, topN, vectorScoreAttribute)

	cmd := client.B().FtSearch().Index(indexName).Query(searchQuery).Sortby(vectorScoreAttribute).
		Params().Nargs(2).NameValue().NameValue("v", string(blob)).
		Dialect(2).Build()

	_, docs, err := client.Do(context.Background(), cmd).AsFtSearch()

	if err != nil {
		log.Fatal("ft.search error - ", err)
	}

	fmt.Println("here are some recommendations for", queryString, "movies:\n======")

	for _, doc := range docs {

		movieDataInJSON := doc.Doc["$"]

		var m Movie
		err = json.Unmarshal([]byte(movieDataInJSON), &m)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("title:", m.Name)
		fmt.Println("description:", m.Desc)

		scoreVal, err := strconv.ParseFloat(doc.Doc[vectorScoreAttribute], 64)
		if err != nil {
			log.Fatal(err)
		}
		score := 1 - scoreVal
		fmt.Println("similarity score:", score)

		fmt.Println("======")
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

type Movie struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}
