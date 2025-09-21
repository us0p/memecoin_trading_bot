package coinprovider

import (
	"net/http"
	"os"

	"memecoin_trading_bot/app/utils"
)

type TokenHolder struct {
	Address        string  `json:"address"`
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

type HeliusResponseList struct {
	Result struct {
		Value []TokenHolder `json:"value"`
	} `json:"result"`
}

func GetTokenLargestHolders(client *http.Client, url, mint string) ([]TokenHolder, error) {
	if mint == "" {
		return []TokenHolder{}, nil
	}
	rpc_params := newHeliusRPCParams(
		"getTokenLargestAccounts",
		mint,
	)

	requester, err := utils.NewRequester[HeliusResponseList](client, url, http.MethodPost)
	if err != nil {
		return []TokenHolder{}, err
	}

	requester.AddHeader("ContentType", "application/json")
	requester.AddQuery("api-key", os.Getenv("HELIUS_API_KEY"))
	req_with_body, err := requester.AddBody(rpc_params)

	if err != nil {
		return []TokenHolder{}, err
	}

	token_authorities, err := req_with_body.Do()
	if err != nil {
		return []TokenHolder{}, err
	}

	return token_authorities.Result.Value, nil
}
