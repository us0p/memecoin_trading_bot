package entities

import "time"

type Trade struct {
	Mint                          string
	IssuedTradeStartAt            time.Time
	TradeStartedAt                time.Time
	IssuedTradeEndAt              time.Time
	TradeEndedAt                  time.Time
	IssuedTradeStartTokenUsdPrice float64
	IssuedTradeEndTokenUsdPrice   float64
	EntryTokenUsdPrice            float64
	ExitTokenUsdPrice             float64
	SolanaAmount                  float64
	ExecutedSolanaAmount          float64
	TotalFees                     float64
	ExpectedTokenAmount           float64
	ExecutedTokenAmount           float64
}
