package workflows

import (
	"context"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/entities"
	"memecoin_trading_bot/app/notification"
	"net/http"
	"sync"
)

func GetTradeOpportunityLargestHolders(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
) {
	ctx := context.Background()
	latest_trade_opp, err := db_client.GetLatestTradeOpp(ctx)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.DatabaseOp,
			err,
			notification.Fatal,
		)
		return
	}

	tokens_top_holders := make([]entities.TopHolder, len(latest_trade_opp))
	var wg sync.WaitGroup
	for idx, mint := range latest_trade_opp {
		wg.Add(1)
		go func(mint string, idx int, db_client *db.DB) {
			defer wg.Done()
			token_largest_holders, err := coinprovider.GetTokenLargestHolders(
				http_client,
				constants.HELIUS_API_URL,
				mint,
			)

			if err != nil {
				nf_state.RecordError(
					mint,
					notification.LargestHolders,
					err,
					notification.Transient,
				)
				return
			}

			ctx := context.Background()

			total_supply, err := db_client.GetLatestSupplyForToken(
				ctx,
				mint,
			)
			if err != nil {
				nf_state.RecordError(
					mint,
					notification.DatabaseOp,
					err,
					notification.Fatal,
				)
				return
			}

			top_5_wallets := aggTopWalletPercentage(
				total_supply,
				token_largest_holders[0:5],
			)
			top_10_wallets := aggTopWalletPercentage(
				total_supply,
				token_largest_holders[0:10],
			)
			top_20_wallets := aggTopWalletPercentage(
				total_supply,
				token_largest_holders[0:20],
			)

			hslh := has_single_largest_holder(
				total_supply,
				token_largest_holders[0].UiAmount,
			)

			tokens_top_holders[idx] = entities.NewTopHolder(
				mint,
				top_5_wallets,
				top_10_wallets,
				top_20_wallets,
				hslh,
			)
		}(mint, idx, db_client)
	}

	wg.Wait()

	if err = db_client.InsertTopHolderBulk(ctx, tokens_top_holders); err != nil {
		nf_state.RecordError(
			"",
			notification.DatabaseOp,
			err,
			notification.Fatal,
		)
		return
	}
}

func aggTopWalletPercentage(
	total_supply float64,
	token_holders []coinprovider.TokenHolder,
) float64 {
	if total_supply == 0 {
		return 0
	}

	var total float64
	for _, holder := range token_holders {
		total += holder.UiAmount
	}

	return (total / total_supply) * 100
}

func has_single_largest_holder(total_supply, biggest_holder_amount float64) bool {
	return biggest_holder_amount > (total_supply / 2)
}
