package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var file string
	flag.StringVar(&file, "file", "", "takes the contents of the file to upload")

	var u string
	flag.StringVar(&u, "u", "", "url to upload the file or the contents of the file")

	var server string
	flag.StringVar(&server, "server", "http://localhost:8080", "server to upload the files")
	flag.Parse()

	if u == "" {
		flag.Usage()
		return
	}

	var content string
	switch file {
	case "":
		c, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error while reading stdin: %v", err)
		}
		content = string(c)
	default:
		c, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Error while reading file: %v", err)
		}
		content = string(c)
	}
	fmt.Println("Content: ", content)

	req, err := http.NewRequest(http.MethodPost, server+"/"+u, strings.NewReader(content))
	if err != nil {
		log.Fatalf("Error sending post request to the server: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error executing port request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("✓ Pasted successfully!\nURL: %s/%s\n", server, u)
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to paste: %s (status %d)", string(body), resp.StatusCode)
	}
}
