package entities

import "time"

type Operation string

const (
	BUY  Operation = "buy"
	SELL Operation = "sell"
)

type Order struct {
	Mint string
	Op   Operation
}

type Trade struct {
	Mint                         string
	Operation                    Operation
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
