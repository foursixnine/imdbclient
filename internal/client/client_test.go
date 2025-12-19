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
	empty_url := imdbClient.makeUrl("", []QueryParameters{})
	if empty_url != url.String() {
		t.Errorf("TestIMDBClientMakeURL() = got (%v), want (%v).", empty_url, url.String())
	}

	params := []QueryParameters{
		{Key: "", Value: ""},
		{Key: "key", Value: "value"},
	}
	non_empty_query := imdbClient.makeUrl("", params)
	expected := url.String() + "?key=value"
	if non_empty_query != expected {
		t.Errorf("TestIMDBClientMakeURL() = got (%v), want (%v).", non_empty_query, expected)
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
