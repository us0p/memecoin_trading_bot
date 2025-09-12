package coinprovider

import (
	"encoding/json"
	"net/http"
)

const memescan_api_url = "https://memescan.app/api/calls?sort=recent&limit=10&offset=0&type=gamble"

type Call struct {
	Mint      string `json:"tokenAddress"`
	Symbol    string `json:"tokenSymbol"`
	CreatedAt string `json:"createdAt"`
}

type MemeScanResponse struct {
	Calls []Call `json:"calls"`
}

func GetGambleTokens(client *http.Client) ([]Call, error) {
	res, err := client.Get(memescan_api_url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var memescan_response MemeScanResponse
	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&memescan_response)

	if err != nil {
		return nil, err
	}

	return memescan_response.Calls, nil
}
