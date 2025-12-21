package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/foursixnine/imdblookup/internal/client"
	"github.com/foursixnine/imdblookup/models"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	fmt.Println("Application started")

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})
	var titles []*models.ImdbapiTitle
	var err error

	url, err := url.Parse("https://api.imdbapi.dev")
	if err != nil {
		panic(err)
	}

	options := client.ImdbClientOptions{
		ApiURL:    url,
		Verbose:   true,
		UserAgent: "imdblookup/0.1",
	}

	imdbClient := client.New(options)
	go func() {
		defer wg.Done()
		fmt.Println("Finding results:")
		titles, err = imdbClient.FindShowsByTitle()
		if err != nil {
			panic(err)
		}
		close(done)
		fmt.Println("\nDone fetching results.")
	}()
	counter := 1
	go func() {
		for {
			select {
			case <-time.After(1 * time.Millisecond):
				fmt.Print(".")
				counter++
				if counter%100 == 0 {
					fmt.Print("\n")
					counter = 0
				}
			case <-done:
				return
			}
		}
	}()

	wg.Wait()

	if err != nil {
		panic(err)
	}

	if len(titles) == 0 {
		log.Println("No titles found")
	}

	for _, title := range titles {
		fmt.Printf("(%s)\t-> \"%s\"\n", title.ID, title.OriginalTitle)
	}

}
