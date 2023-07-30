package client

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {
	client = http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	os.Exit(m.Run())
}

func TestClient_200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	actual, err := Get(ts.URL, Cache{})
	if err != nil {
		log.Fatal(err)
	}

	expected := &Response{StatusCode: 200, Content: "Hello, client\n"}

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("value is mismatch (-actual +expected):\n%s", diff)
	}
}

func TestClient_302(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Add("Location", fmt.Sprintf("http://%s/%s", r.Host, "next"))
			w.WriteHeader(302)
			return
		}
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	actual, err := Get(ts.URL, Cache{})
	if err != nil {
		log.Fatal(err)
	}

	expected := &Response{StatusCode: 200, Content: "Hello, client\n"}

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("value is mismatch (-actual +expected):\n%s", diff)
	}
}

func TestClient_CheckRedirectLoop(t *testing.T) {
	n := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/n") {
			n++
			w.Header().Add("Location", fmt.Sprintf("http://%s/n%d", r.Host, n))
			w.WriteHeader(302)
		}
	}))
	defer ts.Close()

	actual, err := Get(ts.URL, Cache{})
	if err != nil {
		log.Fatal(err)
	}

	expected := &Response{StatusCode: 400}
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("value is mismatch (-actual +expected):\n%s", diff)
	}

	if n != 6 {
		t.Error("redirect loop check counter is not invalid.")
	}
}

func TestClient_304(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") == "etag" && r.Header.Get("If-Modified-Since") == "modified" {
			w.Header().Set("Etag", "etag2")
			w.Header().Set("Last-Modified", "Thu, 14 Jan 2021 12:54:14 GMT")
			fmt.Fprintf(w, "ok")
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	actual, err := Get(ts.URL, Cache{Etag: "etag", Modified: "modified"})
	if err != nil {
		log.Fatal(err)
	}

	expected := &Response{StatusCode: 200, Content: "ok", Cache: Cache{Modified: "Thu, 14 Jan 2021 12:54:14 GMT", Etag: "etag2"}}
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("value is mismatch (-actual +expected):\n%s", diff)
	}
}

func TestClient_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer ts.Close()

	actual, err := Get(ts.URL, Cache{})
	if err != nil {
		log.Fatal(err)
	}

	expected := &Response{StatusCode: 404}
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("value is mismatch (-actual +expected):\n%s", diff)
	}
}

func TestClient_301(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Add("Location", fmt.Sprintf("http://%s/n1", r.Host))
			w.WriteHeader(301)
			return
		}
		if r.URL.Path == "/n1" {
			fmt.Fprintf(w, "ok")
			return
		}
		w.WriteHeader(404)
	}))
	defer ts.Close()

	actual, err := Get(ts.URL, Cache{})
	if err != nil {
		log.Fatal(err)
	}

	expected := &Response{StatusCode: 200, Location: ts.URL + "/n1", Content: "ok"}
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("value is mismatch (-actual +expected):\n%s", diff)
	}
}

func TestClient_301_invalid(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Add("Location", "http://domain.invalid")
			w.WriteHeader(301)
			return
		}
		fmt.Fprintf(w, "ok")
	}))
	defer ts.Close()

	actual, err := Get(ts.URL, Cache{})
	if err != nil {
		log.Fatal(err)
	}

	expected := &Response{StatusCode: 400}
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("value is mismatch (-actual +expected):\n%s", diff)
	}
}
