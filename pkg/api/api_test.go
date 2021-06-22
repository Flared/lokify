package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/flared/lokify/pkg/api"
	"github.com/flared/lokify/pkg/api/appctx"
	"github.com/flared/lokify/pkg/loki"
)

type lokiClientTest struct {
	queryReturn *loki.QueryResponse
}

func (client lokiClientTest) Query(query string) (*loki.QueryResponse, error) {
	return client.queryReturn, nil
}

func newLokiClientTest(queryReturn *loki.QueryResponse) *lokiClientTest {
	return &lokiClientTest{
		queryReturn: queryReturn,
	}
}

type apiTest struct {
	server *httptest.Server
	t      *testing.T
}

func (apiTest *apiTest) Close() {
	apiTest.server.Close()
}

func newApiTest(t *testing.T, lokiClient *lokiClientTest) *apiTest {
	ctx := appctx.New(lokiClient)
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

	lokiClient := newLokiClientTest(expected)

	apiTest := newApiTest(t, lokiClient)
	defer apiTest.Close()

	queryUrl := strings.Join([]string{apiTest.server.URL, "api/query?query=%7Blabel%3D%22value%22%7D"}, "/")

	if resp, err := http.Get(queryUrl); err != nil {
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
