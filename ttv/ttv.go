// Same as ttv.sh but in Go
// Use through ttv wrapper
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func readConfig(path string) []string {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(f), "\n")
}

func main() {
	const (
		TWITCH_URL       = "https://www.twitch.tv/"
		DEFAULT_CFG_PATH = "$HOME/.config/ttv/ttv.rc"
	)
	channels := readConfig(os.ExpandEnv(DEFAULT_CFG_PATH))

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, channel := range channels {
		go func(channel string) {
			resp, err := http.Get(TWITCH_URL + channel)
			defer wg.Done()
			defer resp.Body.Close()
			if err != nil {
				log.Println(err)
			}
			page, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			if strings.Contains(string(page), "isLiveBroadcast") {
				fmt.Println(channel)
			}
		}(channel)
	}
	wg.Wait()
}
