package client

import (
	"encoding/json"

	ce "github.com/foursixnine/imdblookup/internal/errors"
	"github.com/foursixnine/imdblookup/models"
)

func (imdbClient *ImdbClient) FindShowsByTitle(searchTitle string) ([]*models.ImdbapiTitle, *ce.IMDBClientApplicationError) {
	if searchTitle == "" {
		err := ce.NewIMDBClientApplicationError("Search title cannot be empty", nil)
		return nil, err
	}

	var titlesResults models.ImdbapiSearchTitlesResponse
	var titles []*models.ImdbapiTitle

	path := "search/titles"
	parameters := []QueryParameters{
		{Key: "query", Value: searchTitle},
		{Key: "limit", Value: "5"},
	}

	resp, err := imdbClient.Get(path, &parameters)
	if err != nil {
		clientErr, ok := err.(*ce.IMDBClientError)
		if !ok {
			return nil, ce.NewIMDBClientApplicationError("unexpected error type: %w", err)
		}
		return nil, ce.NewIMDBClientApplicationError("An error occurred querying search results", clientErr)
	}

	if err := json.Unmarshal(resp, &titlesResults); err != nil {
		return nil, ce.NewIMDBClientApplicationError("error: JSON answer cannot be read", err)
	}

	titles = titlesResults.Titles
	return titles, nil
}
