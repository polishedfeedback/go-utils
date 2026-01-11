package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var categoryMap = map[string][]string{
	"images":   {".img", ".jpeg", ".jpg", ".png", ".webp", ".ai"},
	"docs":     {".pdf", ".doc", ".docx", ".xlsx", ".xls", ".txt", ".pptx"},
	"videos":   {".mp4", ".avi", ".mkv", ".mov"},
	"audio":    {".mp3", ".wav", ".flac"},
	"archives": {".zip", ".rar", ".7z", ".tar", ".gz"},
	"code":     {".go", ".py", ".js", ".html", ".css"},
}

var Path string

func main() {
	flag.StringVar(&Path, "path", "", "path to organize the files")
	flag.Parse()

	files, err := os.ReadDir(Path)
	if err != nil {
		log.Fatalf("Error while reading the directory: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		src := filepath.Join(Path, file.Name())
		ext := filepath.Ext(src)
		destDir := getDestinationDir(ext)

		if destDir != "" {
			moveFile(src, destDir, file.Name())
		} else {
			log.Printf("Unrecognized file path: %s", ext)
		}
	}
}

func getDestinationDir(extension string) string {
	for dir, exts := range categoryMap {
		for _, ext := range exts {
			if ext == extension {
				return filepath.Join(Path, dir)
			}
		}
	}
	return ""
}

func moveFile(src, destDir, fileName string) {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.Mkdir(destDir, os.ModePerm); err != nil {
			log.Printf("Error creating directory %s: %v", destDir, err)
			return
		}
	}
	destPath := filepath.Join(destDir, fileName)
	err := os.Rename(src, destPath)
	if err != nil {
		log.Printf("Error moving the file %s to %s: %v", fileName, destPath, err)
	}
}
