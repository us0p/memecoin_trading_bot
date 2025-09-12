package coinprovider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jupiter_ultra_api_url = "https://lite-api.jup.ag/ultra/v1"

type MarketData struct {
	Mint        string  `json:"id"`
	Twitter     string  `json:"twitter"`
	Website     string  `json:"website"`
	Telegram    string  `json:"telegram"`
	PriceUsd    float64 `json:"usdPrice"`
	Liquidity   float64 `json:"liquidity"`
	HolderCount int     `json:"holderCount"`
	MarketCap   float64 `json:"mcap"`
	Stats1h     struct {
		VolumeChange      float64 `json:"volumeChange"`
		BuyVolume         float64 `json:"buyVolume"`
		SellVolume        float64 `json:"sellVolume"`
		BuyOrganicVolume  float64 `json:"buyOrganicVolume"`
		SellOrganicVolume float64 `json:"sellOrganicVolume"`
		NumBuys           int     `json:"numBuys"`
		NumSells          int     `json:"numSells"`
		NumTraders        int     `json:"numTraders"`
		NumOrganicBuyers  int     `json:"numOrganicBuyers"`
		NumNetBuyers      int     `json:"numNetBuyers"`
	} `json:"stats1h"`
}

func GetMarketDataForAddresses(client *http.Client, addresses []string) ([]MarketData, error) {
	joined_addresses := strings.Join(addresses, ",")

	url := fmt.Sprintf("%s/search?query=%s", jupiter_ultra_api_url, joined_addresses)

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		var res_err map[string]any

		if err = decoder.Decode(&res_err); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(
			"Received status code %d at '/search' in Jupiter. Error: %s",
			res.StatusCode,
			err,
		)
	}

	var tokens_market_data []MarketData

	if err = decoder.Decode(&tokens_market_data); err != nil {
		return nil, err
	}

	return tokens_market_data, nil
}
