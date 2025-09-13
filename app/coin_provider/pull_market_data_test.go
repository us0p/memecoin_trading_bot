package coinprovider

import (
	"errors"
	"fmt"
	"memecoin_trading_bot/app/app_errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestReturnEarlyOnEmptyAddresses(t *testing.T) {
	has_hit_server := false
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			has_hit_server = true
			w.WriteHeader(200)
		}),
	)

	client := ts.Client()
	addresses := []string{}
	GetMarketDataForAddresses(
		client,
		ts.URL,
		addresses,
	)

	if has_hit_server {
		t.Error("Shouldn't have called server since there's no address.")
	}
}

func TestReturErrOnAPIRequest(t *testing.T) {
	tests := []struct {
		name             string
		responseStatus   int
		err              error
		responseBody     string
		expectMarketData []MarketData
	}{
		{"apiError", http.StatusBadRequest, app_errors.ErrNonOkStatus, `{"error": "very strange error"}`, []MarketData{}},
		{
			"apiSuccess",
			http.StatusOK,
			nil,
			`[
				{
					"id": "address1", 
					"twitter": null,
					"website": null,
					"telegram": null,
					"usdPrice": 5.01,
					"liquidity": 1.005,
					"holderCount": 500,
					"mcap": 100.50,
					"stats1h": {
						"volumeChange": 0.5,
						"buyVolume": 0.5,
						"sellVolume": 0.5,
						"buyOrganicVolume": 0.5,
						"sellOrganicVolume": 0.5,
						"numBuys": 5,
						"numSells": 5,
						"numTraders": 5,
						"numOrganicBuyers": 5,
						"numNetBuyers": 5
					}
				},
				{
					"id": "address2", 
					"twitter": null,
					"website": null,
					"telegram": null,
					"usdPrice": 5.01,
					"liquidity": 1.005,
					"holderCount": 500,
					"mcap": 100.50,
					"stats1h": {
						"volumeChange": 0.5,
						"buyVolume": 0.5,
						"sellVolume": 0.5,
						"buyOrganicVolume": 0.5,
						"sellOrganicVolume": 0.5,
						"numBuys": 5,
						"numSells": 5,
						"numTraders": 5,
						"numOrganicBuyers": 5,
						"numNetBuyers": 5
					}
				}
			]`,
			[]MarketData{
				{
					Mint:        "address1",
					Twitter:     "",
					Website:     "",
					Telegram:    "",
					PriceUsd:    5.01,
					Liquidity:   1.005,
					HolderCount: 500,
					MarketCap:   100.50,
					Stats1h: struct {
						VolumeChange      float64 "json:\"volumeChange\""
						BuyVolume         float64 "json:\"buyVolume\""
						SellVolume        float64 "json:\"sellVolume\""
						BuyOrganicVolume  float64 "json:\"buyOrganicVolume\""
						SellOrganicVolume float64 "json:\"sellOrganicVolume\""
						NumBuys           int     "json:\"numBuys\""
						NumSells          int     "json:\"numSells\""
						NumTraders        int     "json:\"numTraders\""
						NumOrganicBuyers  int     "json:\"numOrganicBuyers\""
						NumNetBuyers      int     "json:\"numNetBuyers\""
					}{
						VolumeChange:      0.5,
						BuyVolume:         0.5,
						SellVolume:        0.5,
						BuyOrganicVolume:  0.5,
						SellOrganicVolume: 0.5,
						NumBuys:           5,
						NumSells:          5,
						NumTraders:        5,
						NumOrganicBuyers:  5,
						NumNetBuyers:      5,
					},
				},
				{
					Mint:        "address2",
					Twitter:     "",
					Website:     "",
					Telegram:    "",
					PriceUsd:    5.01,
					Liquidity:   1.005,
					HolderCount: 500,
					MarketCap:   100.50,
					Stats1h: struct {
						VolumeChange      float64 "json:\"volumeChange\""
						BuyVolume         float64 "json:\"buyVolume\""
						SellVolume        float64 "json:\"sellVolume\""
						BuyOrganicVolume  float64 "json:\"buyOrganicVolume\""
						SellOrganicVolume float64 "json:\"sellOrganicVolume\""
						NumBuys           int     "json:\"numBuys\""
						NumSells          int     "json:\"numSells\""
						NumTraders        int     "json:\"numTraders\""
						NumOrganicBuyers  int     "json:\"numOrganicBuyers\""
						NumNetBuyers      int     "json:\"numNetBuyers\""
					}{
						VolumeChange:      0.5,
						BuyVolume:         0.5,
						SellVolume:        0.5,
						BuyOrganicVolume:  0.5,
						SellOrganicVolume: 0.5,
						NumBuys:           5,
						NumSells:          5,
						NumTraders:        5,
						NumOrganicBuyers:  5,
						NumNetBuyers:      5,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(test.responseStatus)
						fmt.Fprint(w, test.responseBody)
					},
				),
			)

			client := ts.Client()
			addresses := []string{"address1", "address2"}
			mk_data, err := GetMarketDataForAddresses(
				client,
				ts.URL,
				addresses,
			)

			if !errors.Is(err, test.err) {
				t.Errorf("Expected error to be %s, received %s", test.err, err)
			}

			if !slices.Equal(mk_data, test.expectMarketData) {
				t.Errorf(
					"Expected market data to be %+v, received: %+v",
					test.expectMarketData,
					mk_data,
				)
			}
		})
	}
}
