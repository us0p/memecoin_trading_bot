package coinprovider

import (
	"encoding/json"
	"fmt"
	"io"
	"memecoin_trading_bot/app/app_errors"
	"net/http"
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
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		err_resp, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: %s",
				app_errors.ErrReadingAPIResponse,
				err,
			)
		}

		return nil, fmt.Errorf(
			"%w, Received status: %d from MemeScan API. %s",
			app_errors.ErrNonOkStatus,
			res.StatusCode,
			err_resp,
		)
	}

	var memescan_response MemeScanResponse

	err = decoder.Decode(&memescan_response)

	if err != nil {
		return nil, err
	}

	return memescan_response.Calls, nil
}
