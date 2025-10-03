package entities

import "time"

type Order string

const (
	BUY  Order = "buy"
	SELL Order = "sell"
)

type Trade struct {
	Mint                         string
	operation                    Order
	SlippageBPS                  int
	InputAmountLamports          int
	ExpectedOutputAmountLamports int
	InputUSDPrice                float64
	TotalFeeLamports             int
	ExpectedTokenUSDPrice        float64
	IssuedOrderAt                time.Time
	ReceivedOrderResponseAt      time.Time
	ExecutedOutputAmountLamports int
	ExecutedTokenUSDPrice        float64
}
