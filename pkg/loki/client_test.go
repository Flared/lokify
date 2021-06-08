package loki_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/flared/lokify/pkg/loki"
)

func TestQuery(t *testing.T) {
	respBody := &loki.QueryResponse{
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
	var respBodyBytes bytes.Buffer
	if err := json.NewEncoder(&respBodyBytes).Encode(respBody); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expectedUrl := "/loki/api/v1/query?limit=100&query=%7Blabel%3D%22value%22%7D"
		if req.URL.String() != expectedUrl {
			t.Fatalf("got %v, expect %v", req.URL.String(), expectedUrl)
		}
		rw.Write(respBodyBytes.Bytes())
	}))
	defer server.Close()

	lokiClient := loki.NewClient(server.Client(), server.URL)
	if resp, err := lokiClient.Query(`{label="value"}`); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(resp, respBody) {
		t.Fatalf("got %+v, expect %+v", resp, respBody)
	}
}
