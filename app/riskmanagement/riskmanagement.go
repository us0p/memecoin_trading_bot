package riskmanagement

import (
	"context"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/utils"
	"net/http"
)

const (
	MAX_TRADE_PERCENTAGE  = 0.1
	WALLET_FEE_PERCENTAGE = 0.01
)

func GetTradeAmount(http_client *http.Client, db_client *db.DB) (float64, error) {
	pvk, err := utils.GetPrvKey()
	if err != nil {
		return 0, err
	}

	wallet_balance, err := coinprovider.GetOnChainWalletHoldings(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		pvk.PublicKey().String(),
	)
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	ongoing_trades_balance_lamports, err := db_client.GetOngoingTradesBalanceLamports(ctx)
	if err != nil {
		return 0, err
	}

	total_balance := wallet_balance.UiAmount + utils.FromLamports(ongoing_trades_balance_lamports)

	trade_amount := total_balance * MAX_TRADE_PERCENTAGE
	fee_pool := total_balance * WALLET_FEE_PERCENTAGE

	total_available_balance := wallet_balance.UiAmount - fee_pool

	if trade_amount > total_available_balance {
		return total_available_balance, nil
	}

	return trade_amount, nil
}
