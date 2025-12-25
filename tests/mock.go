package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/foursixnine/imdblookup/models"
)

func SetupServer(t *testing.T) (server *httptest.Server) {
	t.Helper()
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range r.Header {
			log.Printf("Request header: %s => %v\n", key, value)
		}
		log.Printf("Request path: %s", r.URL.Path)
		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			fmt.Fprint(w, "Hello, world")
		case "":
			fmt.Fprint(w, "Hello, world")
		case "/foo":
			http.NotFoundHandler().ServeHTTP(w, r)
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal server error")
		case "/search/titles":
			params := r.URL.Query()
			query := params.Get("query")
			data, err := getDataForQuery(t, query)

			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Accept", "application/json")
			w.Header().Add("Accept-Charset", "UTF-8")

			if err != nil {
				fmt.Fprint(w, err)
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			w.Write(data)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}

	}))

	return
}

func getDataForQuery(t *testing.T, query string) (data []byte, err error) {
	t.Helper()
	switch query {
	case "Stranger Things":
		titles := models.ImdbapiSearchTitlesResponse{}
		titleValues := []*models.ImdbapiTitle{
			{ID: "foobar", OriginalTitle: "Stranger Things"},
		}
		titles.Titles = titleValues //append(titles.Titles, &title)
		data, err = json.Marshal(titles)
	case "Broken Json":
		data = []byte("{")
		err = nil
	default:
		data = []byte("{}")
		err = nil
	}
	return
}
