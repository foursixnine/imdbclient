package client

import (
	"encoding/json"
	e "errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/foursixnine/imdblookup/internal/errors"
	"github.com/foursixnine/imdblookup/models"
)

var server *httptest.Server

func TestParametersForQuery(t *testing.T) {
	// Test creating a ParametersForQuery instance
	param := QueryParameters{
		Key:   "query",
		Value: "Stranger Things",
	}

	if param.Key != "query" {
		t.Errorf("Expected Key to be 'query', got '%s'", param.Key)
	}

	if param.Value != "Stranger Things" {
		t.Errorf("Expected Value to be 'query', got '%s'", param.Value)
	}
}

func TestClient(t *testing.T) {
	// url, err := url.Parse(server.URL)
	url, err := url.Parse(server.URL)
	if err != nil {
		t.Error("Failed to parse url")
	}

	options := ImdbClientOptions{
		ApiURL:    url,
		Verbose:   true,
		UserAgent: "imdblookup/0.1",
	}

	imdbClient := New(options)
	if imdbClient.options.ApiURL != url {
		t.Error("Error creating instance of IMDBClient")
	}

	if server.URL != url.String() {
		t.Error("Error creating instance of IMDBClient")
	}

}

func TestIMDBClientGet(t *testing.T) {
	testCases := map[string]struct {
		params   []QueryParameters
		path     string
		expected string
		error    *errors.HTTPError
	}{
		"with empty path, with query": {
			path:     "",
			expected: "Hello, world",
			params: []QueryParameters{
				{Key: "", Value: ""},
				{Key: "key", Value: "value"},
			},
		},
		"with real path, with query, 404": {
			path:     "/foo",
			expected: "Hello, world",
			params: []QueryParameters{
				{Key: "", Value: ""},
				{Key: "key", Value: "value"},
			},
			error: errors.NotFound(server.URL + "foo?key=value"),
		},
		"with real path, with query, 500": {
			path:     "/500",
			expected: "Internal Server Error",
			params: []QueryParameters{
				{Key: "", Value: ""},
				{Key: "key", Value: "value"},
			},
			error: errors.UnexpectedError(http.StatusInternalServerError, server.URL+"foo?key=value"),
		},
	}

	url, err := url.Parse(server.URL)
	if err != nil {
		log.Println("Failed to parse url", err)
	}

	options := ImdbClientOptions{
		ApiURL:    url,
		Verbose:   true,
		UserAgent: "imdblookup/0.1",
	}

	imdbClient := New(options)
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			resp, err := imdbClient.Get(testCase.path, &testCase.params)

			if err != nil && testCase.error == nil {
				t.Fatalf("TestIMDBClientGet(%s) in executing get request: err: (%v) resp:(%v)", testName, err, string(resp))
			} else if err != nil && testCase.error != nil {

				t.Log("Testing errors")
				t.Logf("TestIMDBClientGet(%s) = got (%v), want (%v).", testName, err, testCase.error)

				if e.Is(err, testCase.error) && testCase.error != err {
					t.Errorf("TestIMDBClientGet(%s) = got (%v), want (%v).", testName, err, testCase.error)
				}

			} else if string(resp) != testCase.expected {
				t.Errorf("TestIMDBClientGet(%s) = got (%v), want (%v).", testName, string(resp), testCase.expected)
			}
		})
	}

}

func TestIMDBClientMakeURL(t *testing.T) {
	testCases := map[string]struct {
		params   []QueryParameters
		path     string
		expected string
	}{
		"with empty path, with query": {
			path:     "",
			expected: server.URL + "?key=value",
			params: []QueryParameters{
				{Key: "", Value: ""},
				{Key: "key", Value: "value"},
			},
		},
		"with empty path, no query": {
			path:     "",
			expected: server.URL,
			params: []QueryParameters{
				{Key: "", Value: ""},
			},
		},
		"with path, with query": {
			path:     "/foo",
			expected: server.URL + "/foo?bar=baz&key=value",
			params: []QueryParameters{
				{Key: "bar", Value: "baz"},
				{Key: "key", Value: "value"},
			},
		},
		"with emoji, with query": {
			path:     "/🥺",
			expected: server.URL + "/%F0%9F%A5%BA?bar=baz&key=value",
			params: []QueryParameters{
				{Key: "bar", Value: "baz"},
				{Key: "key", Value: "value"},
			},
		},
		"with emoji, with emoji in query": {
			path:     "/🥺",
			expected: server.URL + "/%F0%9F%A5%BA?bar=baz&key=%F0%9F%A5%BA",
			params: []QueryParameters{
				{Key: "bar", Value: "baz"},
				{Key: "key", Value: "🥺"},
			},
		},
		"with path, with real value in query": {
			path:     "/foo",
			expected: server.URL + "/foo?query=Stranger+Things",
			params: []QueryParameters{
				{Key: "query", Value: "Stranger Things"},
			},
		},
	}

	url, err := url.Parse(server.URL)
	if err != nil {
		t.Errorf("Failed to parse url (%v)", err)
	}

	options := ImdbClientOptions{
		ApiURL:    url,
		Verbose:   true,
		UserAgent: "imdblookup/0.1",
	}

	imdbClient := New(options)
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			result := imdbClient.makeUrl(testCase.path, testCase.params)
			if result != testCase.expected {
				t.Errorf("TestIMDBClientMakeURL(%s) = got (%v), want (%v).", testName, result, testCase.expected)
			}
		})
	}
}

func TestIMDBClientFindShowsByTitle(t *testing.T) {
	testCases := map[string]struct {
		params   string
		expected []*models.ImdbapiTitle
		error    *errors.HTTPError
	}{
		"with empty query": {
			expected: []*models.ImdbapiTitle{},
			params:   "",
		}, "With non empty query": {
			expected: []*models.ImdbapiTitle{
				{ID: "foobar", OriginalTitle: "Stranger Things"},
			},
			params: "Stranger Things",
		},
	}
	url, err := url.Parse(server.URL)
	if err != nil {
		log.Println("Failed to parse url", err)
	}

	options := ImdbClientOptions{
		ApiURL:    url,
		Verbose:   true,
		UserAgent: "imdblookup/0.1",
	}

	imdbClient := New(options)
	for testName, testCase := range testCases {
		titles, err := imdbClient.FindShowsByTitle(testCase.params)
		if err != nil {
			t.Errorf("TestIMDBClientFindShowsByTitle(%v) = got error (%v)", testName, err)
		}

		if len(titles) != len(testCase.expected) {
			t.Fatalf("TestIMDBClientFindShowsByTitle(%v) = Got (%v) more results than expected (%v)", testName, len(titles), len(testCase.expected))
		}
	}

}

func TestMain(m *testing.M) {

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
			titles := models.ImdbapiSearchTitlesResponse{}
			params := r.URL.Query()
			query := params.Get("query")
			if query == "Stranger Things" {
				titleValues := []*models.ImdbapiTitle{
					{ID: "foobar", OriginalTitle: "Stranger Things"},
				}
				titles.Titles = titleValues //append(titles.Titles, &title)
			}
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Accept", "application/json")
			w.Header().Add("Accept-Charset", "UTF-8")
			data, err := json.Marshal(titles)

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
	defer server.Close()

	log.Println("Server started: ", server.URL)

	os.Exit(m.Run())

}
