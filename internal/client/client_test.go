package client

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
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

func TestIMDBClientGet_empty_query(t *testing.T) {
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
	resp, err := imdbClient.Get("", &[]QueryParameters{})
	if err != nil {
		t.Fatalf("error in executing get request: %v", err)
	}
	expected := "Hello, world"
	if string(resp) != expected {
		t.Logf("TestApiServer_empty_query() = got (%v), want (%v).", string(resp), expected)
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

func TestMain(m *testing.M) {

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range r.Header {
			log.Printf("Request header: %s => %v\n", key, value)
		}
		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			// fmt.Println(r.Header)
			fmt.Fprint(w, "Hello, world")
		case "":
			// fmt.Println(r.Header)
			fmt.Fprint(w, "Hello, world")
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}

	}))
	defer server.Close()

	log.Println("Server started: ", server.URL)

	os.Exit(m.Run())

}
