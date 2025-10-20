package workflows

import (
	"context"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/entities"
	"memecoin_trading_bot/app/notification"
	"memecoin_trading_bot/app/riskmanagement"
	"net/http"
)

func GetTradeOpportunityMarketData(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
	orders_chan chan<- entities.Order,
) {
	ctx := context.Background()
	latest_trade_opp, err := db_client.GetLatestTradeOpp(ctx)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.PullMarketData,
			err,
			notification.Fatal,
		)
		return
	}

	mk_data, err := coinprovider.GetMarketDataForAddresses(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		latest_trade_opp,
	)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.PullMarketData,
			err,
			notification.Fatal,
		)
		return
	}

	orders, err := riskmanagement.CheckTradesToClose(db_client, mk_data)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.PullMarketData,
			err,
			notification.Fatal,
		)
		return
	}
	for _, order := range orders {
		err := db_client.InsertTradeProcessing(ctx, order)
		if err != nil {
			nf_state.RecordError(
				order.Mint,
				notification.PullMarketData,
				err,
				notification.Fatal,
			)
			return
		}
		orders_chan <- order
	}

	if err = db_client.InsertMarketDataBulk(ctx, mk_data); err != nil {
		nf_state.RecordError(
			"",
			notification.PullMarketData,
			err,
			notification.Fatal,
		)
		return
	}
}
