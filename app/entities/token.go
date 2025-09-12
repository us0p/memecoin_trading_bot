package entities

import "time"

type Token struct {
	Mint          string    // populated by MemeScan
	Symbol        string    // populated by MemeScan
	CreatedAt     time.Time // populated by MemeScan
	MintEnabled   bool      // populated by Helius getAccountInfo
	FreezeEnabled bool      // populated by Helius getAccountInfo
	TradeOpp      bool      // populated by Jupiter market data
	Twitter       string    // populated by Jupiter market data
	Site          string    // populated by Jupiter market data
	Telegram      string    // populated by Jupiter market data
}
