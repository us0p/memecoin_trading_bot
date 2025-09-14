package coinprovider

import (
	"net/http"

	"memecoin_trading_bot/app/coin_provider/utils"
)

type Call struct {
	Mint      string `json:"tokenAddress"`
	Symbol    string `json:"tokenSymbol"`
	CreatedAt string `json:"createdAt"`
}

type MemeScanResponse struct {
	Calls []Call `json:"calls"`
}

func GetGambleTokens(client *http.Client, url string) ([]Call, error) {
	requester, err := utils.NewRequester[MemeScanResponse](client, url, http.MethodGet)
	if err != nil {
		return nil, err
	}

	memescan_response, err := requester.Do()
	if err != nil {
		return nil, err
	}

	return memescan_response.Calls, nil
}
