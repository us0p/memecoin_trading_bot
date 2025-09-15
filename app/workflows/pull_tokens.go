package workflows

import (
	"context"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/entities"
	"net/http"
	"slices"
)

func validateTradeOpportunity(mkData coinprovider.MarketData) bool {
	if mkData.HolderCount >= 1000 && mkData.Liquidity >= 500000 {
		return true
	}
	return false
}

func PullTokens(http_client *http.Client, db_client *db.DB) {
	ctx := context.Background()

	calls, err := coinprovider.GetGambleTokens(http_client, constants.MEMESCAN_API_URL)
	if err != nil {
	}

	mints := make([]string, len(calls))
	for idx, call := range calls {
		mints[idx] = call.Mint
	}

	registered_tokens, err := db_client.QueryExistingTokensFromSlice(ctx, mints)
	if err != nil {
	}

	newTokens := make([]string, len(mints)-len(registered_tokens))
	for _, mint := range mints {
		if !slices.Contains(registered_tokens, mint) {
			newTokens = append(newTokens, mint)
		}
	}

	mk_data, err := coinprovider.GetMarketDataForAddresses(http_client, constants.JUPITER_ULTRA_API_URL, newTokens)
	if err != nil {
	}

	token_authorities := make([]coinprovider.TokenAuthorities, len(newTokens))
	for _, mint := range newTokens {
		go func() {
			authorities, err := coinprovider.GetTokenAuthorities(http_client, constants.HELIUS_API_URL, mint)
			if err != nil {
			}

			token_authorities = append(token_authorities, authorities)
		}()
	}

	for _, newToken := range newTokens {
		token := entities.Token{
			Mint: newToken,
		}
	}

}
