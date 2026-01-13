package api

import (
	"io"
	"net/http"
	"net/url"
)

const (
	BaseURL   = "https://api.allanime.day/api"
	Referer   = "https://allmanga.to"
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
)

// MakeRequest creates and executes an HTTP request with proper headers
func MakeRequest(gqlQuery, variables string) ([]byte, error) {
	params := url.Values{}
	params.Add("query", gqlQuery)
	params.Add("variables", variables)

	req, err := http.NewRequest("GET", BaseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", Referer)
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
