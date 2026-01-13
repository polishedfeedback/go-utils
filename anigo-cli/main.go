package main

import (
	"fmt"

	"github.com/polishedfeedback/go-utils/anigo-cli/internal/api"
	"github.com/polishedfeedback/go-utils/anigo-cli/internal/decoder"
	"github.com/polishedfeedback/go-utils/anigo-cli/internal/player"
)

func main() {
	fmt.Print("Search anime: ")
	var query string
	fmt.Scanln(&query)

	results, err := api.SearchAnime(query)
	if err != nil || len(results) == 0 {
		fmt.Println("No results found")
		return
	}

	fmt.Println("\nSearch Results:")
	for i, anime := range results {
		fmt.Printf("%d. %s (%d episodes)\n", i+1, anime.Name, anime.Episodes)
	}

	fmt.Print("\nSelect anime number: ")
	var choice int
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(results) {
		fmt.Println("Invalid choice")
		return
	}

	selectedAnime := results[choice-1]
	fmt.Printf("\nSelected: %s\n", selectedAnime.Name)

	episodes, err := api.GetEpisodes(selectedAnime.ID)
	if err != nil || len(episodes) == 0 {
		fmt.Println("No episodes found")
		return
	}

	fmt.Printf("\nAvailable Episodes (%d total):\n", len(episodes))
	for i, ep := range episodes {
		fmt.Printf("%d. Episode %s\n", i+1, ep)
	}

	fmt.Print("\nSelect episode number: ")
	var epChoice int
	fmt.Scanln(&epChoice)

	if epChoice < 1 || epChoice > len(episodes) {
		fmt.Println("Invalid episode")
		return
	}

	selectedEpisode := episodes[epChoice-1]
	fmt.Printf("\nYou selected Episode %s\n", selectedEpisode)

	encodedURL, err := api.GetVideoSources(selectedAnime.ID, selectedEpisode)
	if err != nil {
		fmt.Printf("Error getting video: %v\n", err)
		return
	}

	videoURL := decoder.DecodeURL(encodedURL)
	fmt.Printf("\n✓ Video URL obtained\n")

	title := fmt.Sprintf("%s - Episode %s", selectedAnime.Name, selectedEpisode)
	if err := player.Play(videoURL, title); err != nil {
		fmt.Printf("Error playing video: %v\n", err)
	}
}
