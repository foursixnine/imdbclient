package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/foursixnine/imdblookup/internal/client"
	ce "github.com/foursixnine/imdblookup/internal/errors"
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

	if args.api == "" {
		log.Fatalf("api url cannot be empty")
	} else if !strings.HasPrefix(args.api, "http") {
		log.Fatalf("api url does not have scheme: '%s'", args.api)
	}

	url, err := url.Parse(args.api)
	if err != nil {
		log.Fatalf("Error parsing api url: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})
	result := &ce.IMDBClientApplicationError{}

	imdbClient := client.New(url)
	log.Printf("Application started with %s as Server\n", args.api)

	go getTitles(imdbClient, opts.Query, &wg, done, result)
	go progressMarker(done)

	wg.Wait()

	if result.Code != 0 {
		log.Printf("Error not empty, %v\n", result)
		ce.RootCause(result)
		os.Exit(result.Code)
	}
}
