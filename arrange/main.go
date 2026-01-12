package main

import (
	"flag"
	"fmt"
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

var (
	Path    string
	Preview bool
)

func main() {
	flag.StringVar(&Path, "path", "", "path to organize the files")
	flag.BoolVar(&Preview, "preview", false, "preview changes without executing")
	flag.Parse()

	if Path == "" {
		log.Fatal("Error: --path is required")
	}
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		log.Fatalf("Error: Path '%s' doesn't exist", Path)
	}

	files, err := os.ReadDir(Path)
	if err != nil {
		log.Fatalf("Error while reading the directory: %v", err)
		return
	}

	fileCount := 0
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
		}
	}
	if fileCount == 0 {
		fmt.Println("No files to organize")
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
			if Preview {
				previewMove(destDir, file.Name())
				fmt.Println("\n--- Preview Mode ---")
				fmt.Println("Run without --preview flag to execute these changes")
			} else {
				moveFile(src, destDir, file.Name())
			}
		} else {
			if Preview {
				previewMove(filepath.Join(Path, "others"), file.Name())
				fmt.Println("\n--- Preview Mode ---")
				fmt.Println("Run without --preview flag to execute these changes")
			} else {
				moveFile(src, filepath.Join(Path, "others"), file.Name())
			}
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

	if _, err := os.Stat(destPath); err == nil {
		destPath = getUniqueFilename(destPath)
	}

	err := os.Rename(src, destPath)
	if err != nil {
		log.Printf("Error moving file %s: %v", fileName, err)
		return
	}

	fmt.Printf("✓ Moved %s → %s\n", fileName, destPath)
}

func getUniqueFilename(path string) string {
	ext := filepath.Ext(path)
	nameWithoutExt := path[:len(path)-len(ext)]

	counter := 1
	newPath := path

	for {
		newPath = fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break
		}
		counter++
	}

	return newPath
}

func previewMove(destDir, fileName string) {
	destPath := filepath.Join(destDir, fileName)

	if _, err := os.Stat(destPath); err == nil {
		uniquePath := getUniqueFilename(destPath)
		fmt.Printf("[PREVIEW] %s → %s (renamed to avoid conflict)\n", fileName, uniquePath)
	} else {
		fmt.Printf("[PREVIEW] %s → %s\n", fileName, destPath)
	}
}
