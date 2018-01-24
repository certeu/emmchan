// https://github.com/google/go-github/blob/master/github/github_test.go
package emm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the EMM client being tested.
	client *Client

	// server is a test HTTP server used to provide mock RSS responses.
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// client to use test server
	client = NewClient(nil)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v; want %v", r.Method, m)
		}
		fmt.Fprint(w, `response body`)
	})
	resp, err := client.Get(fmt.Sprintf("%s/feed", server.URL))
	if err != nil {
		t.Fatal(err)
	}
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != "response body" {
		t.Errorf("Response body = %v; want response body", string(actual))
	}
}
