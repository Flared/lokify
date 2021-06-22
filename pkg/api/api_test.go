package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/flared/lokify/pkg/api"
	"github.com/flared/lokify/pkg/loki"
)

type lokiClientTest struct {
	output     *loki.QueryResponse
	testInputs func(...string)
}

func (client lokiClientTest) Query(query string) (*loki.QueryResponse, error) {
	client.testInputs(query)
	return client.output, nil
}

func (client lokiClientTest) QueryRange(query string, start string, end string) (*loki.QueryResponse, error) {
	client.testInputs(query, start, end)
	return client.output, nil
}

type apiTest struct {
	server *httptest.Server
	t      *testing.T
}

func (apiTest *apiTest) Close() {
	apiTest.server.Close()
}

func newApiTest(t *testing.T, lokiClient *lokiClientTest) *apiTest {
	ctx := api.NewAppContext(lokiClient)
	return &apiTest{
		server: httptest.NewServer(api.NewRouter(ctx)),
		t:      t,
	}
}

func TestApi(t *testing.T) {
	apiTest := newApiTest(t, nil)
	defer apiTest.Close()

	if _, err := http.Get(strings.Join([]string{apiTest.server.URL, "api/status"}, "/")); err != nil {
		t.Fatal(err)
	}

}

func TestQuery(t *testing.T) {
	expected := &loki.QueryResponse{
		Status: "success",
		Data: loki.QueryData{
			ResultType: "streams",
			Result: []loki.QueryDataResult{
				{
					Stream: map[string]string{"label": "value"},
					Values: [][2]string{{"123456789", "log message"}},
				},
			},
		},
	}

	lokiClient := &lokiClientTest{
		output: expected,
		testInputs: func(values ...string) {
			if !contains(values, `{label="value"}`) {
				t.Fatal("query is missing")
			}
		},
	}

	apiTest := newApiTest(t, lokiClient)
	defer apiTest.Close()

	if resp, err := http.Get(strings.Join([]string{apiTest.server.URL, "api/query?query=%7Blabel%3D%22value%22%7D"}, "/")); err != nil {
		t.Fatal(err)
	} else {
		defer resp.Body.Close()

		var q loki.QueryResponse
		if err := json.NewDecoder(resp.Body).Decode(&q); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(&q, expected) {
			t.Fatalf("got %+v, expect %+v", &q, expected)
		}
	}
}

func contains(values []string, lookup string) bool {
	for _, value := range values {
		if value == lookup {
			return true
		}
	}
	return false
}

func TestQueryRange(t *testing.T) {
	expected := &loki.QueryResponse{
		Status: "success",
		Data: loki.QueryData{
			ResultType: "streams",
			Result: []loki.QueryDataResult{
				{
					Stream: map[string]string{"label": "value"},
					Values: [][2]string{{"123456789", "log message"}},
				},
			},
		},
	}

	lokiClient := &lokiClientTest{
		output: expected,
		testInputs: func(values ...string) {
			if !contains(values, `{label="value"}`) {
				t.Fatal("query is missing")
			}
			if !contains(values, "123456789") {
				t.Fatal("start is missing")
			}
			if !contains(values, "123456678") {
				t.Fatal("end is missing")
			}
		},
	}

	apiTest := newApiTest(t, lokiClient)
	defer apiTest.Close()

	queryRangeUrl := strings.Join(
		[]string{apiTest.server.URL, "api/query_range?query=%7Blabel%3D%22value%22%7D&start=123456789&end=123456678"},
		"/",
	)
	if resp, err := http.Get(queryRangeUrl); err != nil {
		t.Fatal(err)
	} else {
		defer resp.Body.Close()

		var q loki.QueryResponse
		if err := json.NewDecoder(resp.Body).Decode(&q); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(&q, expected) {
			t.Fatalf("got %+v, expect %+v", &q, expected)
		}
	}
}
