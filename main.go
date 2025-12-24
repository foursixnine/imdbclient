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
	var err *ce.IMDBClientApplicationError

	url, parseerr := url.Parse("https://api.imdbapi.dev")
	if parseerr != nil {
		panic(parseerr)
	}

	options := client.ImdbClientOptions{
		ApiURL:    url,
		Verbose:   true,
		UserAgent: "imdblookup/0.1",
	}

	imdbClient := client.New(options)
	go func() {
		var titles []*models.ImdbapiTitle

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
	}()
	go func() {
		progressMarker(done)
	}()

	wg.Wait()

	if err != nil {
		panic(err)
	}

}

func progressMarker(done chan struct{}) {
	counter := 0
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
}
