package api

import (
	"encoding/json"
	"fmt"

	"github.com/polishedfeedback/go-utils/anigo-cli/models"
)

// SearchAnime searches for anime by query
func SearchAnime(query string) ([]models.Anime, error) {
	gqlQuery := `query($search:SearchInput $limit:Int $page:Int $translationType:VaildTranslationTypeEnumType $countryOrigin:VaildCountryOriginEnumType){shows(search:$search limit:$limit page:$page translationType:$translationType countryOrigin:$countryOrigin){edges{_id name availableEpisodes __typename}}}`
	variables := fmt.Sprintf(`{"search":{"allowAdult":false,"allowUnknown":false,"query":"%s"},"limit":40,"page":1,"translationType":"sub","countryOrigin":"ALL"}`, query)

	body, err := MakeRequest(gqlQuery, variables)
	if err != nil {
		return nil, err
	}

	var result models.SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var animeList []models.Anime
	for _, anime := range result.Data.Shows.Edges {
		animeList = append(animeList, models.Anime{
			ID:       anime.ID,
			Name:     anime.Name,
			Episodes: anime.Episodes.Sub,
		})
	}

	return animeList, nil
}

// GetEpisodes fetches episode list for an anime
func GetEpisodes(animeID string) ([]string, error) {
	gqlQuery := `query($showId:String!){show(_id:$showId){_id availableEpisodesDetail}}`
	variables := fmt.Sprintf(`{"showId":"%s"}`, animeID)

	body, err := MakeRequest(gqlQuery, variables)
	if err != nil {
		return nil, err
	}

	var result models.EpisodesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	episodes := result.Data.Show.AvailableEpisodesDetail.Sub

	for i, j := 0, len(episodes)-1; i < j; i, j = i+1, j-1 {
		episodes[i], episodes[j] = episodes[j], episodes[i]
	}

	return episodes, nil
}

// GetVideoSources fetches video sources for an episode
func GetVideoSources(animeID, episodeNum string) (string, error) {
	gqlQuery := `query($showId:String!,$translationType:VaildTranslationTypeEnumType!,$episodeString:String!){episode(showId:$showId translationType:$translationType episodeString:$episodeString){episodeString sourceUrls}}`
	variables := fmt.Sprintf(`{"showId":"%s","translationType":"sub","episodeString":"%s"}`, animeID, episodeNum)

	body, err := MakeRequest(gqlQuery, variables)
	if err != nil {
		return "", err
	}

	var result models.EpisodeSourceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	for _, source := range result.Data.Episode.SourceUrls {
		if source.Type == "player" {
			return source.SourceUrl, nil
		}
	}

	if len(result.Data.Episode.SourceUrls) > 0 {
		return result.Data.Episode.SourceUrls[0].SourceUrl, nil
	}

	return "", fmt.Errorf("no video sources found")
}
