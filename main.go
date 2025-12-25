package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/foursixnine/imdblookup/internal/client"
	ce "github.com/foursixnine/imdblookup/internal/errors"
	"github.com/foursixnine/imdblookup/models"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.Println("Application started")

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})

	url, err := url.Parse("https://api.imdbapi.dev")
	if err != nil {
		log.Printf("Error parsing api url: %v", err)
		os.Exit(1)
	}

	options := client.ImdbClientOptions{
		ApiURL:    url,
		Verbose:   true,
		UserAgent: "imdblookup/0.1",
	}

	imdbClient := client.New(options)
	go getTitles(imdbClient, &wg, done)
	go progressMarker(done)

	wg.Wait()

	// if err != nil {
	// 	panic(err)
	// }

}

func getTitles(imdbClient *client.ImdbClient, wg *sync.WaitGroup, done chan struct{}) {
	var titles []*models.ImdbapiTitle
	var err *ce.IMDBClientApplicationError

	defer wg.Done()
	fmt.Println("Finding results:")
	titles, err = imdbClient.FindShowsByTitle("Stranger Things")
	if err != nil {
		var appErr *ce.IMDBClientApplicationError
		if errors.As(err, &appErr) && err.AppMessage == "Search title cannot be empty" {
			log.Printf("Title cannot be empty %v\n", err)
			os.Exit(404)
		}
		log.Printf("An error has occurred: (%v)\n", err)
		os.Exit(2)
	}
	close(done)
	fmt.Println("\nDone fetching results.")

	if len(titles) == 0 {
		log.Println("No titles found")
	}

	for _, title := range titles {
		fmt.Printf("(%s)\t-> \"%s\"\n", title.ID, title.OriginalTitle)
	}
}

func progressMarker(done chan struct{}) {
	counter := 0
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
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
}
