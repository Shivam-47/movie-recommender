package main

import (
	"fmt"
	"os"

	pgdump "github.com/Shivam-47/movie-recommender/postgres/dump"
	pgrecommend "github.com/Shivam-47/movie-recommender/postgres/recommend"
	rdump "github.com/Shivam-47/movie-recommender/redis/dump"
	rindex "github.com/Shivam-47/movie-recommender/redis/index"
	rrecommend "github.com/Shivam-47/movie-recommender/redis/recommend"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

func main() {
	var mode string
	var query string
	var topN int

	rootCmd := &cobra.Command{
		Use:   "movie-recommender",
		Short: "Movie recommender CLI",
		Run: func(cmd *cobra.Command, args []string) {
			switch mode {
			case "postgres-dump":
				pgdump.Main()
			case "postgres-recommend":
				pgrecommend.Main(query, topN)
			case "redis-init":
				rindex.Main()
			case "redis-dump":
				rdump.Main()
			case "redis-recommend":
				rrecommend.Main(query, topN)
			default:
				fmt.Println("unknown mode:", mode)
				fmt.Println("use --mode postgres-dump, postgres-recommend, redis-init, redis-dump, or redis-recommend")
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVar(&mode, "mode", "", "mode to run: postgres-dump | postgres-recommend | redis-init | redis-dump | redis-recommend")
	rootCmd.Flags().StringVar(&query, "query", "bruce wayne", "search text used by the recommend modes")
	rootCmd.Flags().IntVar(&topN, "top", 5, "number of recommendations returned by the recommend modes")
	rootCmd.MarkFlagRequired("mode")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
