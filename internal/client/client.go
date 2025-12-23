package client

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/foursixnine/imdblookup/internal/errors"
)

type QueryParameters struct {
	Key   string
	Value string
}

type ImdbClientOptions struct {
	ApiURL    *url.URL
	Verbose   bool
	UserAgent string
}

type ImdbClient struct {
	HttpClient *http.Client
	options    *ImdbClientOptions
}

type imdbClientTransport struct {
	UserAgent string
}

func New(options ImdbClientOptions) *ImdbClient {
	transport := &imdbClientTransport{
		UserAgent: options.UserAgent,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	return &ImdbClient{
		HttpClient: httpClient,
		options:    &options,
	}
}

func (t *imdbClientTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Accept-Charset", "UTF-8")
	r.Header.Add("User-Agent", t.UserAgent)
	// r.Header.Add("X-AUTH-API-KEY", t.apiKey)

	return http.DefaultTransport.RoundTrip(r)
}

func (client *ImdbClient) Get(path string, params *[]QueryParameters) ([]byte, error) {

	url := client.makeUrl(path, *params)
	log.Println("ImdbClient querying: " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.NewIMDBClientGenericError("error: creating http request %w", err)
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, errors.NewIMDBClientGenericError("error: executing http request %w", err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewIMDBClientGenericError("error: reading body of request %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return response, errors.NotFound(url)
		} else {
			return response, errors.UnexpectedError(resp.StatusCode, string(response))
		}
	}

	return response, nil
}

func (client *ImdbClient) makeUrl(path string, params []QueryParameters) string {
	url := client.options.ApiURL.JoinPath(path)
	q := url.Query()
	for _, query := range params {
		if query.Key != "" {
			q.Set(query.Key, query.Value) // Set the key-value in the URL's query parameters
		}
	}
	url.RawQuery = q.Encode()
	return url.String()
}
