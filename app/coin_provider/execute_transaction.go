package coinprovider

import (
	"memecoin_trading_bot/app/utils"
	"net/http"
)

type JupiterOrderExecutionResp struct {
	Status            string `json:"status"`
	Signature         string `json:"signature"`
	InputAmountResult string `json:"inputAmountResult"`
	OutputAmountResut string `json:"outputAmountResult"`
}

type orderExecutionParams struct {
	SignedTransaction string `json:"signedTransaction"`
	RequestId         string `json:"requestId"`
}

func ExecuteTransaction(
	http_client *http.Client,
	url,
	signedTransaction,
	requestId string,
) (JupiterOrderExecutionResp, error) {
	requester, err := utils.NewRequester[JupiterOrderExecutionResp](http_client, url, http.MethodPost)
	if err != nil {
		return JupiterOrderExecutionResp{}, err
	}

	requester.AddPath("/execute")
	requester.AddHeader("Content-Type", "application/json")
	requester.AddBody(orderExecutionParams{signedTransaction, requestId})

	resp, err := requester.Do()
	return resp, err
}
