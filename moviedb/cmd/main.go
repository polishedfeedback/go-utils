package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/polishedfeedback/moviedb/internal/storage"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  movies add --year <year> --actors <actor1,actor2> <movie title>")
	fmt.Println("  movies list")
	fmt.Println("  movies get <id>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  movies add --year 2010 --actors \"DiCaprio,Hardy,Page\" \"Inception\"")
	fmt.Println("  movies list")
	fmt.Println("  movies get 1")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	store, err := storage.NewSQLiteStorage("movies.db")
	if err != nil {
		fmt.Errorf("Error connecting to the database: %v", err)
		return
	}

	switch command {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		year := addCmd.Int("year", 0, "Movie Year")
		actorsStr := addCmd.String("actors", "", "Comma-Separated actors")

		addCmd.Parse(os.Args[2:])

		if addCmd.NArg() < 1 {
			fmt.Println("Error: title required")
			printUsage()
			return
		}

		title := addCmd.Arg(0)
		actors := strings.Split(*actorsStr, ",")
		fmt.Printf("DEBUG: actorsStr = '%s'\n", *actorsStr)
		fmt.Printf("DEBUG: actors slice = %v\n", actors)
		fmt.Printf("DEBUG: number of actors = %d\n", len(actors))

		err := store.AddMovie(title, *year, actors)
		if err != nil {
			log.Fatalf("Error adding movie: %v", err)
		}

		fmt.Printf("✓ Added movie: %s (%d)\n", title, *year)

	case "get":
		getCmd := flag.NewFlagSet("get", flag.ExitOnError)
		getCmd.Parse(os.Args[2:])

		if getCmd.NArg() < 1 {
			fmt.Println("Error: id required")
			return
		}

		id, err := strconv.Atoi(getCmd.Arg(0))
		if err != nil {
			log.Fatalf("error converting id: %v", err)
		}

		movie, err := store.GetMovie(id)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("\nMovie: %s (%d)\n", movie.Title, movie.Year)
		fmt.Printf("Actors: %s\n", strings.Join(movie.Actors, ", "))

	case "list":
		movies, err := store.ListMovies()
		if err != nil {
			log.Fatalf("error getting movie list: %v", err)
		}
		for _, movie := range movies {
			fmt.Printf("\nMovie: %s (%d)\n", movie.Title, movie.Year)
			fmt.Printf("Actors: %s\n", strings.Join(movie.Actors, ", "))
		}
	}
}
