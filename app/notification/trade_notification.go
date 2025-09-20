package notification

import (
	"fmt"
	"memecoin_trading_bot/app/constants"
	"time"
)

type TradeOpening struct {
	Symbol           string
	OpenedAt         time.Time
	SolAmount        float64
	WalletPercentage float64
}

type Status string

const (
	Win  Status = "WIN"
	Loss Status = "LOSS"
)

type TradeClosing struct {
	Symbol           string
	ClosedAt         time.Time
	Duration         string
	ProfitPercentage float32
	ProfitSOL        float64
	OpStatus         Status
}

func (n *Notifications) openTradeReport() []string {
	reports := make([]string, len(n.TradesOpened))

	for idx, t := range n.TradesOpened {
		reports[idx] = fmt.Sprintf(`*TRADE OPENING*
			- *Symbol*: %s
			- *Opened at*: %s
			- *Amount SOL*: %f
			- *Wallet percentage*: %f%%`,
			t.Symbol,
			t.OpenedAt.Format(constants.NOTIFICATION_TIME_REP),
			t.SolAmount,
			t.WalletPercentage,
		)
	}

	return reports
}

func (n *Notifications) closeTradeReport() []string {
	reports := make([]string, len(n.TradesClosed))

	for idx, t := range n.TradesClosed {
		reports[idx] = fmt.Sprintf(`*TRADE CLOSED*
			- *Symbol*: %s
			- *Closed at*: %s
			- *Duration*: %s
			- *Profit %%*: %f%%
			- *Pofit in SOL*: %f
			- *Operation Status:* %s`,
			t.Symbol,
			t.ClosedAt.Format(constants.NOTIFICATION_TIME_REP),
			t.Duration,
			t.ProfitPercentage,
			t.ProfitSOL,
			t.OpStatus,
		)
	}

	return reports
}
