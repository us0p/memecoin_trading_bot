package coinprovider

import (
	"fmt"
	"memecoin_trading_bot/app/utils"
	"net/http"
	"strconv"
)

type SwapInfo struct {
	AmmKey     string `json:"ammKey"`
	Label      string `json:"label"`
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	InAmount   string `json:"inAmount"`
	OutAmount  string `json:"outAmount"`
	FeeAmount  string `json:"feeAmount"`
	FeeMint    string `json:"feeMint"`
}

type RoutePlan struct {
	SwapInf SwapInfo `json:"swapInfo"`
}

type JupiterGetOrderResponse struct {
	InAmount             string      `json:"inAmount"`
	OutAmount            string      `json:"outAmount"`
	OtherAmountThreshold string      `json:"otherAmountThreshold"`
	SlippageBps          int         `json:"slippageBps"`
	InUsdValue           float64     `json:"inUsdValue"`
	RoutePln             []RoutePlan `json:"routePlan"`
	Transaction          string      `json:"transaction"`
	RequestId            string      `json:"requestId"`
	TotalTime            int         `json:"totalTime"`
	ExpireAt             string      `json:"expireAt"`
	ErrorCode            int         `json:"errorCode"`
	ErrorMessage         string      `json:"errorMessage"`
}

func (j JupiterGetOrderResponse) GetTotalFeeAmount() (int, error) {
	total := 0
	for _, plan := range j.RoutePln {
		as_int, err := strconv.Atoi(plan.SwapInf.FeeAmount)
		if err != nil {
			return 0, err
		}
		total += as_int
	}

	return total, nil
}

func GetTradeTransaction(
	client *http.Client,
	url,
	taker_addrs,
	input_mint,
	output_mint string,
	input_amount int,
) (JupiterGetOrderResponse, error) {
	requester, err := utils.NewRequester[JupiterGetOrderResponse](
		client,
		url,
		http.MethodGet,
	)
	var resp JupiterGetOrderResponse
	if err != nil {
		return resp, err
	}

	requester.AddPath("/order")
	requester.AddQuery("inputMint", input_mint)
	requester.AddQuery("outputMint", output_mint)
	requester.AddQuery("taker", taker_addrs)
	requester.AddQuery("amount", fmt.Sprint(input_amount))

	token_order_dt, err := requester.Do()
	if err != nil {
		return resp, err
	}

	return token_order_dt, nil
}
