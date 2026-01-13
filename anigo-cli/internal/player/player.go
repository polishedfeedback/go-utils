package player

import (
	"fmt"
	"os"
	"os/exec"
)

const referer = "https://allmanga.to"

// Play plays a video URL using mpv or vlc
func Play(videoURL, title string) error {
	fmt.Printf("\n🎬 Playing: %s\n", title)

	if _, err := exec.LookPath("mpv"); err == nil {
		cmd := exec.Command("mpv",
			"--referrer="+referer,
			"--force-media-title="+title,
			videoURL)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if _, err := exec.LookPath("vlc"); err == nil {
		cmd := exec.Command("vlc",
			"--http-referrer="+referer,
			videoURL)
		return cmd.Run()
	}

	return fmt.Errorf("no video player found. Please install mpv or vlc")
}
