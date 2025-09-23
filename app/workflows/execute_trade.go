package workflows

import (
	"context"
	"fmt"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/notification"
	"memecoin_trading_bot/app/riskmanagement"
	"memecoin_trading_bot/app/utils"
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
		return
	}

	pvk, err := utils.GetPrvKey()
	if err != nil {
		nf_state.RecordError(
			"",
			notification.ExecuteTrade,
			err,
			notification.Fatal,
		)
		return
	}

	transactions := make([]string, len(mints))
	var wg sync.WaitGroup
	for idx, mint := range mints {
		wg.Add(1)
		go func(mint string, http_client *http.Client, db_client *db.DB, tx_idx int) {
			defer wg.Done()
			// call riskmanagement to get total token amount
			// get wallet address
			sol_amount, err := riskmanagement.GetTradeAmount(
				http_client,
				db_client,
			)
			if err != nil {
				nf_state.RecordError(
					mint,
					notification.ExecuteTrade,
					err,
					notification.Fatal,
				)
				return
			}
			agg_resp, err := coinprovider.GetTradeTransaction(
				http_client,
				constants.JUPITER_ULTRA_API_URL,
				pvk.PublicKey().String(),
				mint,
				utils.ToLamports(sol_amount),
			)
			if err != nil {
				nf_state.RecordError(
					mint,
					notification.ExecuteTrade,
					err,
					notification.Core,
				)
				return
			}
			if agg_resp.ErrorCode != 0 {
				nf_state.RecordError(
					mint,
					notification.ExecuteTrade,
					fmt.Errorf("Received the following error for transaction %s.", agg_resp.ErrorMessage),
					notification.Core,
				)
				return
			}
			// records trade info into object to save into DB later.
			// sign obj transaction and adds to transaction queue.
			transactions[idx] = agg_resp.Transaction
		}(mint, http_client, db_client, idx)
	}
	wg.Wait()
}
