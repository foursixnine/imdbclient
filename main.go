package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	client := &http.Client{}

	go func() {
		defer wg.Done()
		fmt.Println("Finding results:")
		titles, err = findResults(client)
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

func findResults(client *http.Client) ([]*models.ImdbapiTitle, error) {
	var titlesResults models.ImdbapiSearchTitlesResponse
	var titles []*models.ImdbapiTitle
	// curl -X 'GET' \
	// 'https://api.imdbapi.dev/search/titles?query=Stranger%20Things' \
	// -H 'accept: application/json'
	//models.ImdbapiSearchTitlesResponse

	// req := http.Request{Method: "GET"}

	// resp, err := http.Get("https://api.imdbapi.dev/search/titles?query=Stranger%20Things")
	req, err := http.NewRequest("GET", "https://api.imdbapi.dev/search/titles?query=Stranger%20Things", nil)
	req.Header.Set("User-Agent", "imdblookup/0.1")

	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		return titles, fmt.Errorf("Error: Status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return titles, fmt.Errorf("Error: Status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&titlesResults); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return titles, errors.New("Document can't be read")
	}
	titles = titlesResults.Titles
	return titles, nil
}
