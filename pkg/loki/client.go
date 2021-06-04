package loki

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

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
	u, errBuildUrl := buildUrl(
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

	resp, errGet := client.client.Get(u.String())
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
