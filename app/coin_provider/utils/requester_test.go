package utils

import (
	"errors"
	"fmt"
	"memecoin_trading_bot/app/app_errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestResponse struct {
	Success string `json:"success"`
}

func TestServerResponseTreatment(t *testing.T) {
	tests := []struct {
		name             string
		statusCode       int
		respBody         string
		expectedRespBody TestResponse
		expectedError    error
	}{
		{"serverError", http.StatusBadRequest, `{"error": "crazy error"}`, TestResponse{}, app_errors.ErrNonOkStatus},
		{"serverSuccess", http.StatusOK, `{"success": "OK"}`, TestResponse{"OK"}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(test.statusCode)
						fmt.Fprint(w, test.respBody)
					},
				),
			)

			client := ts.Client()
			url := ts.URL
			requester, _ := NewRequester[TestResponse](client, url, http.MethodGet)

			test_response, err := requester.Do()
			if !errors.Is(err, test.expectedError) {
				t.Errorf("Expected error: %s, got: %s", test.expectedError, err)
			}

			if test_response != test.expectedRespBody {
				t.Errorf("Expected response to be: %+v, got: %+v", test.expectedRespBody, test_response)
			}

		})
	}
}
