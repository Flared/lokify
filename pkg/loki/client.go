package loki

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LokiClient struct {
	baseUrl string
}

func NewClient(baseUrl string) *LokiClient {
	return &LokiClient{
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
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func (c LokiClient) Query(query string) (*QueryResponse, error) {
	u, errBuildUrl := buildUrl(
		c.baseUrl,
		"loki/api/v1/query",
		map[string]string{
			"query": query,
			"limit": "100",
		},
	)
	if errBuildUrl != nil {
		return nil, errBuildUrl
	}

	resp, errGet := http.Get(u.String())
	if errGet != nil {
		return nil, errGet
	}

	defer resp.Body.Close()

	b, errReadAll := ioutil.ReadAll(resp.Body)
	if errReadAll != nil {
		return nil, errReadAll
	}

	var q QueryResponse
	if err := json.Unmarshal(b, &q); err != nil {
		return nil, err
	}

	return &q, nil
}
