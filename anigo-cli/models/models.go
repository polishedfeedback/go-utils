package models

// For search results
type Anime struct {
	ID       string
	Name     string
	Episodes int
}

// For search responses from the API
type SearchResponse struct {
	Data struct {
		Shows struct {
			Edges []struct {
				ID       string `json:"_id"`
				Name     string `json:"name"`
				Episodes struct {
					Sub int `json:"sub"`
				} `json:"availableEpisodes"`
			} `json:"edges"`
		} `json:"shows"`
	} `json:"data"`
}

// For episode responses from the API
type EpisodesResponse struct {
	Data struct {
		Show struct {
			ID                      string `json:"_id"`
			AvailableEpisodesDetail struct {
				Sub []string `json:"sub"`
				Dub []string `json:"dub"`
			} `json:"availableEpisodesDetail"`
		} `json:"show"`
	} `json:"data"`
}

// For episode source response from the API
type EpisodeSourceResponse struct {
	Data struct {
		Episode struct {
			EpisodeString string `json:"episodeString"`
			SourceUrls    []struct {
				SourceUrl  string `json:"sourceUrl"`
				SourceName string `json:"sourceName"`
				Type       string `json:"type"`
			} `json:"sourceUrls"`
		} `json:"episode"`
	} `json:"data"`
}
