package index

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var c *redis.Client

const defaultRedisUsername = "default"

func Main() {

	redisHostPort := os.Getenv("REDIS_HOST_PORT")
	password := os.Getenv("REDIS_PASSWORD")

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

	err = createIndex()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Index created -", indexName)
}

const indexName = "idx:movies"
const indexPrefix = "movie:"
const dimension = 3072

func createIndex() error {

	// drop any existing index along with its documents (DD); errors if it doesn't exist yet, which is fine
	c.Do(context.Background(), "FT.DROPINDEX", indexName, "DD").Result()

	_, err := c.Do(context.Background(), "FT.CREATE", indexName, "ON", "JSON", "PREFIX", "1", indexPrefix, "SCORE", "1.0", "SCHEMA",
		"$.name", "TEXT", "WEIGHT", 1.0, "NOSTEM", //name field in json
		"$.description", "TEXT", "WEIGHT", 1.0, //description field in json
		//description field's embedding in json
		"$.description_embedding", "AS", "vector", "VECTOR", "FLAT", 6, "TYPE", "FLOAT32", "DIM", dimension, "DISTANCE_METRIC", "COSINE").Result()

	if err != nil {
		return err
	}

	return nil

}
