package coinprovider

import (
	"memecoin_trading_bot/app/utils"
	"net/http"
)

type WalletHoldingToken struct {
	Account        string  `json:"account"`
	Amount         string  `json:"amount"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
	IsFrozen       bool    `json:"isFrozen"`
	Decimals       int     `json:"decimals"`
}

type JupiterOnChainWalletHoldings struct {
	Amount         string                        `json:"amount"`
	UiAmount       float64                       `json:"uiAmount"`
	UiAmountString string                        `json:"uiAmountString"`
	Tokens         map[string]WalletHoldingToken `json:"tokens"`
}

func GetOnChainWalletHoldings(
	client *http.Client,
	url,
	wallet_addrs string,
) (JupiterOnChainWalletHoldings, error) {
	requester, err := utils.NewRequester[JupiterOnChainWalletHoldings](
		client,
		url,
		http.MethodGet,
	)
	if err != nil {
		return JupiterOnChainWalletHoldings{}, err
	}

	requester.AddPath("/holdings")
	requester.AddPath(wallet_addrs)

	wallet_holdings, err := requester.Do()
	if err != nil {
		return JupiterOnChainWalletHoldings{}, err
	}

	return wallet_holdings, nil
}
