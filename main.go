package main

import (
	"encoding/json"
	"fmt"
	"github.com/foursixnine/imdblookup/models"
	"net/http"
)

func main() {
	fmt.Println("Hello world")
	// curl -X 'GET' \
	// 'https://api.imdbapi.dev/search/titles?query=Stranger%20Things' \
	// -H 'accept: application/json'
	//models.ImdbapiSearchTitlesResponse

	resp, err := http.Get("https://api.imdbapi.dev/search/titles?query=Stranger%20Things")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Status code", resp.StatusCode)
		return
	}

	var titles models.ImdbapiSearchTitlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&titles); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _,title := range titles.Titles {
		
		fmt.Printf("(%s)\t-> \"%s\" \n", title.ID, title.OriginalTitle)
		// fmt.Printf("found ", title)

	}

}
