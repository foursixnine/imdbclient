package client

import (
	"encoding/json"
	// "errors"
	"fmt"
	"log"

	"github.com/foursixnine/imdblookup/models"
)

func (imdbClient *ImdbClient) FindShowsByTitle() ([]*models.ImdbapiTitle, error) {
	var titlesResults models.ImdbapiSearchTitlesResponse

	path := "search/titles"
	parameters := []QueryParameters{
		{Key: "query", Value: "Stranger Things"},
		{Key: "limit", Value: "5"},
	}

	resp, err := imdbClient.Get(path, &parameters)
	if err != nil {
		log.Println("An error has occured querying search results")
		panic(err)
	}

	if err := json.Unmarshal(resp, &titlesResults); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, fmt.Errorf("error: Json answer cannot be read: %w", err)
	}

	titles := titlesResults.Titles
	return titles, nil
}
