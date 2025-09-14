package coinprovider

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturnEarlyOnEmptyAddress(t *testing.T) {
	has_hit_server := false
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			has_hit_server = true
			w.WriteHeader(200)
		}),
	)

	client := ts.Client()
	addresses := ""
	GetTokenAuthorities(
		client,
		ts.URL,
		addresses,
	)

	if has_hit_server {
		t.Error("Shouldn't have called server since there's no address.")
	}
}
