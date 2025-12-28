package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/foursixnine/imdblookup/internal/client"
	ce "github.com/foursixnine/imdblookup/internal/errors"
)

func getTitles(imdbClient *client.ImdbClient, query string, wg *sync.WaitGroup, done chan struct{}, result *ce.IMDBClientApplicationError) {
	defer wg.Done()
	fmt.Println("Finding results:")
	titles, err := imdbClient.FindShowsByTitle(&query)

	if err != nil {
		if err.AppMessage == "Search title cannot be empty" {
			log.Printf("Title cannot be empty %v\n", err)
			*result = *err
			result.Code = ce.EMPTYQUERYERROR
			return
		}
		log.Printf("An unexpected error has occurred: (%v)\n", err)
		*result = *err
		result.Code = ce.GENERICERROR
		return
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
