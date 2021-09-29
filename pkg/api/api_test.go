package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/flared/lokify/pkg/api"
)

type lokiQueryResponse struct {
	Status string        `json:"status"`
	Data   lokiQueryData `json:"data"`
}

type lokiQueryData struct {
	ResultType string                `json:"resultType"`
	Result     []lokiQueryDataResult `json:"result"`
}

type lokiQueryDataResult struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

type apiTest struct {
	server *httptest.Server
	t      *testing.T
}

func (apiTest *apiTest) Close() {
	apiTest.server.Close()
}

func newApiTest(t *testing.T, lokiServer *httptest.Server) *apiTest {
	ctx := api.NewContext(lokiServer.Client(), lokiServer.URL)

	return &apiTest{
		server: httptest.NewServer(api.NewRouter(ctx)),
		t:      t,
	}
}

func TestApi(t *testing.T) {
	apiTest := newApiTest(t, &httptest.Server{})
	defer apiTest.Close()

	statusUrl := strings.Join([]string{apiTest.server.URL, "api/status"}, "/")
	if _, err := http.Get(statusUrl); err != nil {
		t.Fatal(err)
	}

}

func TestQuery(t *testing.T) {
	respBody := &lokiQueryResponse{
		Status: "success",
		Data: lokiQueryData{
			ResultType: "streams",
			Result: []lokiQueryDataResult{
				{
					Stream: map[string]string{"label": "value"},
					Values: [][2]string{{"123456789", "log message"}},
				},
			},
		},
	}
	var respBodyBytes bytes.Buffer
	if err := json.NewEncoder(&respBodyBytes).Encode(respBody); err != nil {
		t.Fatal(err)
	}

	lokiServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(respBodyBytes.Bytes())
	}))
	defer lokiServer.Close()

	apiTest := newApiTest(t, lokiServer)
	defer apiTest.Close()

	query := `{label="value"}`

	params := strings.Join(
		[]string{
			"query=" + url.QueryEscape(query),
		},
		"&",
	)

	queryUrl := apiTest.server.URL + "/api/query?" + params

	fmt.Printf("query url, %v\n", queryUrl)

	if resp, err := http.Get(queryUrl); err != nil {
		t.Fatal(err)
	} else {
		defer resp.Body.Close()

		var q lokiQueryResponse
		if err := json.NewDecoder(resp.Body).Decode(&q); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(&q, respBody) {
			t.Fatalf("got %+v, expect %+v", &q, respBody)
		}
	}
}
