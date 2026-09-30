package dump

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/generative-ai-go/genai"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

var movieData map[string]string = map[string]string{
	"The Princess Bride":       "A young woman's fantastical journey to reunite with her true love, facing swashbuckling pirates, a scheming prince, and a fearsome giant.",
	"The Shawshank Redemption": "The story of an innocent man's unwavering hope and resilience as he fights for his freedom and redemption over decades in a brutal prison system.",
	"The Godfather":            "An epic tale of family, loyalty, and betrayal set against the backdrop of the powerful Corleone crime family in post-war New York.",
	"The Dark Knight":          "A gritty and gripping crime thriller that pits Batman against his most formidable foe, the Joker, as they grapple with themes of chaos, order, and justice.",
	"Spirited Away":            "A young girl's journey of self-discovery and courage as she navigates a magical world filled with spirits, witches, and fantastical creatures.",
	"The Lord of the Rings: The Fellowship of the Ring": "A fellowship of hobbits, elves, dwarves, and humans embark on a perilous quest to destroy the One Ring and defeat the evil Sauron.",
	"The Matrix":               "A mind-bending exploration of reality, choice, and destiny as a computer hacker uncovers the truth about the world he lives in.",
	"Inception":                "A thrilling heist adventure that delves into the depths of the human mind and the power of dreams.",
	"Interstellar":             "A crew of astronauts undertakes a daring mission to save humanity by traveling through a wormhole in search of a new home for mankind.",
	"Parasite":                 "A darkly comedic and suspenseful tale of two families from different social classes whose lives become intertwined through manipulation and deception.",
	"The Green Mile":           "A supernatural drama set on death row as a gentle giant with miraculous abilities forms a bond with the guards and inmates.",
	"Forrest Gump":             "The life story of a simple man with extraordinary experiences, touching on themes of love, war, destiny, and the American dream.",
	"The Silence of the Lambs": "A gripping psychological thriller as a young FBI agent seeks the help of an imprisoned cannibalistic serial killer to catch another elusive murderer.",
	"Pulp Fiction":             "A nonlinear narrative of crime, redemption, and dark humor intersecting the lives of hitmen, a boxer, a mob boss's wife, and a pair of petty criminals.",
	"The Usual Suspects":       "A twist-filled crime mystery as a group of criminals come together under mysterious circumstances, leading to a shocking revelation.",
	"The Lion King":            "An animated classic following the journey of a young lion prince reclaiming his rightful place as king after the murder of his father.",
	"The Godfather Part II":    "Continuing the saga of the Corleone family, exploring themes of power, betrayal, and the complexities of the criminal underworld.",
	"The Dark Knight Rises":    "The final chapter of Christopher Nolan's Batman trilogy, depicting the rise of a new enemy and the ultimate test of Gotham's protector.",
	"Fight Club":               "A gritty exploration of masculinity, consumerism, and societal rebellion as an insomniac office worker forms an underground fight club.",
	"Schindler's List":         "A poignant retelling of the true story of Oskar Schindler, a German businessman who saved the lives of Jewish refugees during the Holocaust.",
}

const defaultRedisUsername = "default"

var c *redis.Client
var em *genai.EmbeddingModel

const indexPrefix = "movie:"

func Main() {

	redisHostPort := os.Getenv("REDIS_HOST_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")

	c = redis.NewClient(&redis.Options{
		Addr:     redisHostPort,
		Password: password,
		Username: defaultRedisUsername,
	})

	defer c.Close()

	err := c.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("connect failed: ", err)
	}

	genaiClient, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer genaiClient.Close()

	em = genaiClient.EmbeddingModel("gemini-embedding-001")

	var num int

	for movieName, desc := range movieData {

		num++
		key := indexPrefix + strconv.Itoa(num)

		embedding, err := getEmbedding(desc)
		if err != nil {
			log.Fatal(err)
		}

		movie := Movie{Name: movieName, Desc: desc, Embedding: embedding}
		movieJSON, err := json.Marshal(movie)
		if err != nil {
			log.Fatal(err)
		}

		err = c.JSONSet(context.Background(), key, "$", movieJSON).Err()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("added data with embedding for movie", movieName, "with key", key)
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
	Name      string    `json:"name"`
	Desc      string    `json:"description"`
	Embedding []float32 `json:"description_embedding"`
}
