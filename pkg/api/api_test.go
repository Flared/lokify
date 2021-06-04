package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flared/lokify/pkg/api"
)

type apiTest struct {
	server *httptest.Server
	t      *testing.T
}

func (apiTest *apiTest) Close() {
	apiTest.server.Close()
}

func newApiTest(t *testing.T) *apiTest {
	return &apiTest{
		server: httptest.NewServer(api.NewRouter()),
		t:      t,
	}
}

func TestApi(t *testing.T) {
	apiTest := newApiTest(t)
	defer apiTest.Close()

	if _, err := http.Get(strings.Join([]string{apiTest.server.URL, "api/status"}, "/")); err != nil {
		t.Fatal(err)
	}

}
