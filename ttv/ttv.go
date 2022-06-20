package main
// Same as ttv.sh but in Go*
// Not quite, doesn't read rc file
// Use through ttv wrapper


import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {
	const TWITCH_URL = "https://www.twitch.tv/"

	channel_names := []string{
		"",
	}

	var wg sync.WaitGroup
	for _, channel := range channel_names {
		wg.Add(1)
		go func(channel string) {
			resp, err := http.Get(TWITCH_URL + channel)
			defer resp.Body.Close()
			defer wg.Done()
			if err != nil {
				log.Println(err)
			}
			page, err := ioutil.ReadAll(resp.Body)
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
