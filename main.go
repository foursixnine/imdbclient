package main

import (
	"flag"
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

type CLIargs struct {
	api string
}

type CLIopts struct {
	Query string
	Limit int
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)

	var args CLIargs
	var opts CLIopts

	flag.StringVar(&opts.Query, "query", "Stranger Things", "Search query for IMDB titles")
	flag.StringVar(&args.api, "api", "https://api.imdbapi.dev", "Api url to use as base")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})

	url, err := url.Parse(args.api)
	if err != nil {
		log.Printf("Error parsing api url: %v", err)
		os.Exit(1)
	}

	log.Printf("Application started with %s as Server\n", args.api)
	imdbClient := client.New(url)
	go getTitles(imdbClient, opts.Query, &wg, done)
	go progressMarker(done)

	wg.Wait()

}

func getTitles(imdbClient *client.ImdbClient, query string, wg *sync.WaitGroup, done chan struct{}) {
	var titles []*models.ImdbapiTitle
	var err *ce.IMDBClientApplicationError

	defer wg.Done()
	fmt.Println("Finding results:")
	titles, err = imdbClient.FindShowsByTitle(&query)
	if err != nil {
		if err.AppMessage == "Search title cannot be empty" {
			log.Printf("Title cannot be empty %v\n", err)
			os.Exit(3)
		}
		log.Printf("An error has occurred: (%v)\n", err)
		os.Exit(1)
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
