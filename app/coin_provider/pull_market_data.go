package coinprovider

import (
	"memecoin_trading_bot/app/coin_provider/utils"
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
		return nil, nil
	}

	joined_addresses := strings.Join(addresses, ",")

	requester, err := utils.NewRequester[[]MarketData](client, url, http.MethodGet)
	if err != nil {
		return nil, err
	}

	tokens_market_data, err := requester.AddPath("/search").AddQuery("query", joined_addresses).Do()
	if err != nil {
		return nil, err
	}

	return tokens_market_data, nil
}
