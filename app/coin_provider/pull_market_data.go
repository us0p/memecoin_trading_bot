package coinprovider

import (
	"encoding/json"
	"fmt"
	"io"
	"memecoin_trading_bot/app/app_errors"
	"net/http"
	"strings"
)

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

func GetMarketDataForAddresses(client *http.Client, url string, addresses []string) ([]MarketData, error) {
	if len(addresses) == 0 {
		return []MarketData{}, nil
	}

	joined_addresses := strings.Join(addresses, ",")

	url_with_query := fmt.Sprintf("%s/search?query=%s", url, joined_addresses)

	res, err := client.Get(url_with_query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		resp_body, err := io.ReadAll(res.Body)
		if err != nil {
			return []MarketData{}, fmt.Errorf(
				"%w, %s",
				app_errors.ErrReadingAPIResponse,
				err,
			)
		}

		return nil, fmt.Errorf(
			"%w, Status %d at '/search' from Jupiter. Error: %s",
			app_errors.ErrNonOkStatus,
			res.StatusCode,
			resp_body,
		)
	}

	var tokens_market_data []MarketData

	if err = decoder.Decode(&tokens_market_data); err != nil {
		return nil, err
	}

	return tokens_market_data, nil
}
