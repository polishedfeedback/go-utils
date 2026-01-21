package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	var s string
	flag.StringVar(&s, "s", "", "short url for redirecting")
	var u string
	flag.StringVar(&u, "u", "", "url to use for redirecting")

	var server string
	flag.StringVar(&server, "server", "http://localhost:8080", "server to use")
	flag.Parse()
	if s == "" || u == "" {
		flag.Usage()
		return
	}

	req, err := http.NewRequest(http.MethodPost, server+"/"+s, strings.NewReader(u))
	if err != nil {
		log.Fatalf("error making a post request: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error making a request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Printf("✓ Short URL created successfully!\nURL: %s/%s\n", server, s)
	} else {
		body, _ := io.ReadAll(res.Body)
		log.Fatalf("Failed to create URL: %s (status %d)", string(body), res.StatusCode)
	}
}
