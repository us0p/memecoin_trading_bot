package coinprovider

import (
	"errors"
	"fmt"
	"memecoin_trading_bot/app/app_errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenAuthoritiesApiResponses(t *testing.T) {
	tests := []struct {
		name                string
		statusCode          int
		expectedError       error
		responseBody        string
		expectedAuthorities TokenAuthorities
	}{
		{
			"apiError",
			http.StatusBadRequest,
			app_errors.ErrNonOkStatus,
			`{"error": "some crazy error"}`,
			TokenAuthorities{},
		},
		{
			"apiSuccess",
			http.StatusOK,
			nil,
			`{
				"result": {
					"value": {
						"data": {
							"parsed": {
								"info": {
									"freezeAuthority": null,
									"mintAuthority": null
								}
							}
						}
					}
				}
			}`,
			TokenAuthorities{
				FreezeAuthority: "",
				MintAuthority:   "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(test.statusCode)
						fmt.Fprint(w, test.responseBody)
					},
				),
			)

			client := ts.Client()
			addrs := "mintAddress"

			authorities, err := GetTokenAuthorities(
				client,
				ts.URL,
				addrs,
			)

			if !errors.Is(err, test.expectedError) {
				t.Errorf(
					"Expected error to be %s, got %s",
					test.expectedError,
					err,
				)
			}

			if authorities != test.expectedAuthorities {
				t.Errorf(
					"Expected authorities to be %+v, got %+v",
					test.expectedAuthorities,
					authorities,
				)
			}
		})
	}
}
