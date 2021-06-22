package loki

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Client interface {
	Query(string) (*QueryResponse, error)
}

type lokiClient struct {
	client  *http.Client
	baseUrl string
}

func NewClient(client *http.Client, baseUrl string) *lokiClient {
	return &lokiClient{
		client:  client,
		baseUrl: baseUrl,
	}
}

func buildUrl(baseUrl string, path string, values map[string]string) (*url.URL, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	u.Path += path

	params := url.Values{}

	for key, value := range values {
		params.Add(key, value)
	}

	u.RawQuery = params.Encode()

	return u, nil
}

type QueryResponse struct {
	Status string    `json:"status"`
	Data   QueryData `json:"data"`
}

type QueryData struct {
	ResultType string            `json:"resultType"`
	Result     []QueryDataResult `json:"result"`
}

type QueryDataResult struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

func (client lokiClient) Query(query string) (*QueryResponse, error) {
	queryUrl, errBuildUrl := buildUrl(
		client.baseUrl,
		"loki/api/v1/query",
		map[string]string{
			"query": query,
			"limit": "100",
		},
	)
	if errBuildUrl != nil {
		return nil, errBuildUrl
	}

	resp, errGet := client.client.Get(queryUrl.String())
	if errGet != nil {
		return nil, errGet
	}

	defer resp.Body.Close()

	var q QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&q); err != nil {
		return nil, err
	}

	return &q, nil
}
