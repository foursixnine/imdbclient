package client

import (
	"encoding/json"
	"os"

	// "errors"
	"fmt"
	"log"

	"github.com/foursixnine/imdblookup/models"
)

func (imdbClient *ImdbClient) FindShowsByTitle(searchTitle string) ([]*models.ImdbapiTitle, error) {
	var titlesResults models.ImdbapiSearchTitlesResponse

	path := "search/titles"
	parameters := []QueryParameters{
		{Key: "query", Value: searchTitle},
		{Key: "limit", Value: "5"},
	}

	resp, err := imdbClient.Get(path, &parameters)
	if err != nil {
		log.Printf("An error has occured querying search results (%v)", err)
		os.Exit(2)
	}

	if err := json.Unmarshal(resp, &titlesResults); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, fmt.Errorf("error: Json answer cannot be read: %w", err)
	}

	titles := titlesResults.Titles
	return titles, nil
}
