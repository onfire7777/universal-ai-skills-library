package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClientDropsAPIKeyOnCrossHostRedirect(t *testing.T) {
	seen := make(chan string, 1)
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen <- r.Header.Get("X-API-Key")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer target.Close()

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer source.Close()

	req, err := http.NewRequest(http.MethodGet, source.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-API-Key", "secret-token")

	resp, err := newHTTPClient(5*time.Second, nil).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if got := <-seen; got != "" {
		t.Fatalf("redirect target received X-API-Key %q", got)
	}
}
