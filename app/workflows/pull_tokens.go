package workflows

import (
	"context"
	"fmt"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/entities"
	"memecoin_trading_bot/app/notification"
	"net/http"
	"slices"
	"sync"
	"time"
)

func validateTradeOpportunity(mkData coinprovider.MarketData) bool {
	if mkData.HolderCount >= 1000 && mkData.Liquidity >= 500000 {
		return true
	}
	return false
}

func PullTokens(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
) {
	ctx := context.Background()

	calls, err := coinprovider.GetGambleTokens(http_client, constants.MEMESCAN_API_URL)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.PullCoin,
			err,
			notification.Core,
		)
		return
	}

	mints := make([]string, len(calls))
	for idx, call := range calls {
		mints[idx] = call.Mint
	}

	registered_tokens, err := db_client.QueryExistingTokensFromSlice(ctx, mints)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.PullCoin,
			err,
			notification.Fatal,
		)
		return
	}

	newTokens := make([]coinprovider.Call, len(mints)-len(registered_tokens))
	for idx, call := range calls {
		if !slices.Contains(registered_tokens, call.Mint) {
			newTokens[idx] = call
		}
	}

	mk_data, err := coinprovider.GetMarketDataForAddresses(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		get_token_mints(newTokens),
	)
	if err != nil {
		nf_state.RecordError(
			"",
			notification.TradeOpEval,
			err,
			notification.Core,
		)
		return
	}

	token_authorities := make([]coinprovider.TokenAuthorities, len(newTokens))
	var wg sync.WaitGroup
	for _, token := range newTokens {
		wg.Add(1)
		go func(token coinprovider.Call) {
			defer wg.Done()
			authorities, err := coinprovider.GetTokenAuthorities(
				http_client,
				constants.HELIUS_API_URL,
				token.Mint,
			)
			if err != nil {
				nf_state.RecordError(
					token.Mint,
					notification.TokenAuthorityEval,
					err,
					notification.Transient,
				)
				return
			}

			token_authorities = append(token_authorities, authorities)
		}(token)
	}
	wg.Wait()

	for _, newToken := range newTokens {
		time_rep, err := time.Parse(constants.JAVASCRIPT_TIME_REP, newToken.CreatedAt)
		if err != nil {
			nf_state.RecordError(
				newToken.Mint,
				notification.DatabaseOp,
				err,
				notification.Fatal,
			)
			return
		}

		token_mk_data := get_dt_for_token(
			mk_data,
			newToken.Mint,
		)

		if token_mk_data == nil {
			nf_state.RecordError(
				newToken.Mint,
				notification.TokenDataAgg,
				fmt.Errorf("unmatched mk data for mint: %s", newToken.Mint),
				notification.Fatal,
			)
			return
		}

		token_authority_data := get_dt_for_token(
			token_authorities,
			newToken.Mint,
		)

		if token_authority_data == nil {
			nf_state.RecordError(
				newToken.Mint,
				notification.TokenDataAgg,
				fmt.Errorf("unmatched authority data for mint: %s", newToken.Mint),
				notification.Fatal,
			)
			return
		}

		token := entities.Token{
			Mint:          newToken.Mint,
			Symbol:        newToken.Symbol,
			CreatedAt:     time_rep,
			MintEnabled:   token_authority_data.MintAuthority != "",
			FreezeEnabled: token_authority_data.FreezeAuthority != "",
			TradeOpp:      validateTradeOpportunity(*token_mk_data),
			Twitter:       token_mk_data.Twitter,
			Site:          token_mk_data.Website,
			Telegram:      token_mk_data.Telegram,
		}

		_, err = db_client.InsertToken(ctx, token)
		if err != nil {
			nf_state.RecordError(
				token.Mint,
				notification.DatabaseOp,
				err,
				notification.Fatal,
			)
			return
		}
	}
}

func get_dt_for_token[T coinprovider.HasMint](dt_set []T, mint string) *T {
	for _, dt := range dt_set {
		if dt.GetTokenMint() == mint {
			return &dt
		}
	}

	return nil
}

func get_token_mints(calls []coinprovider.Call) []string {
	mints := make([]string, len(calls))

	for idx, call := range calls {
		mints[idx] = call.Mint
	}

	return mints
}
