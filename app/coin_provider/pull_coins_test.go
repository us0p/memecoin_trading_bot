package coinprovider

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"memecoin_trading_bot/app/app_errors"
)

func TestMemeScanErrors(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		apiStatusCode  int
		calls          []Call
		err            error
	}{
		{
			"apiError",
			`{"error": "crazy error"}`,
			http.StatusBadRequest,
			[]Call{},
			app_errors.ErrNonOkStatus,
		},
		{
			"apiSuccess",
			`{"calls": [{"tokenAddress": "addrs", "tokenSymbol": "symbol", "createdAt": "some date"}]}`,
			http.StatusOK,
			[]Call{
				{"addrs", "symbol", "some date"},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.apiStatusCode)
					fmt.Fprint(w, tt.serverResponse)
				},
			))
			defer ts.Close()

			client := ts.Client()

			calls, err := GetGambleTokens(client, ts.URL)
			if !slices.Equal(calls, tt.calls) {
				t.Errorf(
					"Expected calls to be %+v, received %+v",
					tt.calls,
					calls,
				)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf(
					"Expected error to be %s, received %s",
					tt.err,
					err,
				)
			}
		})
	}

}
