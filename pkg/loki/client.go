package loki

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

	var r QueryResponse
	if errUnmarshal := json.Unmarshal(b, &r); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return &r, nil
}

func (c LokiClient) QueryRange(query string, start string, end string) (*QueryResponse, error) {
	u, errBuildUrl := buildUrl(
		c.baseUrl,
		"loki/api/v1/query_range",
		map[string]string{
			"query": query,
			"limit": "100",
			"start": start,
			"end":   end,
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

	var r QueryResponse
	if errUnmarshal := json.Unmarshal(b, &r); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return &r, nil
}

type LabelsResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

func (c LokiClient) Labels() (*LabelsResponse, error) {
	u, errBuildUrl := buildUrl(
		c.baseUrl,
		"loki/api/v1/labels",
		make(map[string]string),
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

	var r LabelsResponse
	if errUnmarshal := json.Unmarshal(b, &r); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return &r, nil
}

type LabelValuesResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

func (c LokiClient) LabelValues(label string) (*LabelValuesResponse, error) {
	u, errBuildUrl := buildUrl(
		c.baseUrl,
		strings.Join([]string{"loki/api/v1/label", label, "values"}, "/"),
		make(map[string]string),
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

	var r LabelValuesResponse
	if errUnmarshal := json.Unmarshal(b, &r); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return &r, nil
}
