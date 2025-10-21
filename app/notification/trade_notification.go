package notification

import (
	"fmt"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/entities"
	"memecoin_trading_bot/app/utils"
	"time"
)

type TradeOpening struct {
	Symbol           string
	OpenedAt         time.Time
	EntryUSDValue    float64
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
	ProfitPercentage float64
	OpStatus         Status
}

func (n *Notifications) recordTradeOpening(trade_notification_data entities.TradeNotificationData, total_sol_wallet float64) {
	trade_opening_notification := TradeOpening{
		Symbol:           trade_notification_data.Symbol,
		OpenedAt:         trade_notification_data.ReceivedOrderResponseAt,
		EntryUSDValue:    trade_notification_data.InputUSDPrice,
		WalletPercentage: (utils.FromLamports(trade_notification_data.InputAmountLamports) / total_sol_wallet) * 100.0,
	}

	n.TradesOpened = append(n.TradesOpened, trade_opening_notification)
}

func (n *Notifications) recordTradeClosing(trade_buy, trade_sell entities.TradeNotificationData) {
	duration := trade_sell.ReceivedOrderResponseAt.Sub(trade_buy.ReceivedOrderResponseAt)
	profit_percentage := ((trade_sell.ExecutedTokenUSDPrice - trade_buy.ExecutedTokenUSDPrice) / trade_buy.ExecutedTokenUSDPrice) * 100.0
	var op_status Status
	if profit_percentage > 0 {
		op_status = Win
	} else {
		op_status = Loss
	}
	trade_closing_notification := TradeClosing{
		Symbol:           trade_sell.Symbol,
		ClosedAt:         trade_sell.ReceivedOrderResponseAt,
		Duration:         duration.String(),
		ProfitPercentage: profit_percentage,
		OpStatus:         op_status,
	}

	n.TradesClosed = append(n.TradesClosed, trade_closing_notification)
}

func (n *Notifications) openTradeReport() []string {
	reports := make([]string, len(n.TradesOpened))

	for idx, t := range n.TradesOpened {
		reports[idx] = fmt.Sprintf(`*TRADE OPENING*
			- *Symbol*: %s
			- *Opened at*: %s
			- *Entry USD Value*: $%.2f
			- *Wallet percentage*: %.2f%%`,
			t.Symbol,
			t.OpenedAt.Format(constants.NOTIFICATION_TIME_REP),
			t.EntryUSDValue,
			t.WalletPercentage,
		)
	}

	n.TradesOpened = make([]TradeOpening, 0)

	return reports
}

func (n *Notifications) closeTradeReport() []string {
	reports := make([]string, len(n.TradesClosed))

	for idx, t := range n.TradesClosed {
		reports[idx] = fmt.Sprintf(`*TRADE CLOSED*
			- *Symbol*: %s
			- *Closed at*: %s
			- *Duration*: %s
			- *Profit*: %.2f%%
			- *Operation Status:* %s`,
			t.Symbol,
			t.ClosedAt.Format(constants.NOTIFICATION_TIME_REP),
			t.Duration,
			t.ProfitPercentage,
			t.OpStatus,
		)
	}

	(*n).TradesClosed = make([]TradeClosing, 0)

	return reports
}
