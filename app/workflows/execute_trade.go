package workflows

import (
	"context"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/notification"
	"net/http"
	"sync"
)

func ExecuteTrade(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
) {
	ctx := context.Background()
	mints, err := db_client.GetNewTradeMints(ctx)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.ExecuteTrade,
			err,
			notification.Fatal,
		)
	}

	var wg sync.WaitGroup
	for _, mint := range mints {
		wg.Add(1)
		go func(mint string) {
			defer wg.Done()
			// call riskmanagement to get total token amount
			// get wallet address
			agg_resp, err := coinprovider.GetTradeTransaction(
				http_client,
				constants.JUPITER_ULTRA_API_URL,
				"",
				mint,
				0.0,
			)
			if err != nil {
				nf_state.RecordError(
					mint,
					notification.ExecuteTrade,
					err,
					notification.Core,
				)
			}
			// records trade info into object to save into DB later.
			// sign obj transaction and adds to transaction queue.
		}(mint)
	}
	wg.Wait()
}
