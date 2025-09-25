package workflows

import (
	"context"
	"fmt"
	"log"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/entities"
	"memecoin_trading_bot/app/notification"
	"memecoin_trading_bot/app/riskmanagement"
	"memecoin_trading_bot/app/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
)

type tradeOrderCreation struct {
	Trade       entities.Trade
	Transaction string
}

func ExecuteTrade(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
	mint string,
) {
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

	tradeOrderCreation, err := getOrder(
		http_client,
		db_client,
		pvk,
		mint,
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

	if err = signTransaction(pvk, &tradeOrderCreation); err != nil {
		nf_state.RecordError(
			mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}

	ctx := context.Background()
	trade, err := executeOrder(ctx, http_client, db_client, &tradeOrderCreation)
	if err != nil {
		nf_state.RecordError(
			mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}

	if err = db_client.InsertTrade(ctx, trade); err != nil {
		nf_state.RecordError(
			mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}
}

func executeOrder(
	ctx context.Context,
	http_client *http.Client,
	db_client *db.DB,
	tradeOrder *tradeOrderCreation,
) (entities.Trade, error) {
	tradeOrder.Trade.IssuedTradeStartAt = time.Now()

	simu, err := coinprovider.SimulateTransactionExecution(
		http_client,
		constants.HELIUS_API_URL,
		tradeOrder.Transaction,
	)
	if err != nil {
		return entities.Trade{}, err
	}

	tradeOrder.Trade.TradeStartedAt = time.Now()

	last_price, err := db_client.GetLastPriceForToken(ctx, tradeOrder.Trade.Mint)
	if err != nil {
		return entities.Trade{}, err
	}

	tradeOrder.Trade.IssuedTradeStartTokenUsdPrice = last_price

	mk_data, err := coinprovider.GetMarketDataForAddresses(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		[]string{tradeOrder.Trade.Mint},
	)
	if err != nil {
		return entities.Trade{}, err
	}

	tradeOrder.Trade.EntryTokenUsdPrice = mk_data[0].PriceUsd

	logTradeSimulation(simu)

	return tradeOrder.Trade, nil
}

func signTransaction(pvk solana.PrivateKey, tradeOrder *tradeOrderCreation) error {
	tx, err := solana.TransactionFromBase64(tradeOrder.Transaction)
	if err != nil {
		return err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if pvk.PublicKey().Equals(key) {
				return &pvk
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func getOrder(
	http_client *http.Client,
	db_client *db.DB,
	pvk solana.PrivateKey,
	mint string,
) (tradeOrderCreation, error) {
	sol_amount, err := riskmanagement.GetTradeAmount(
		http_client,
		db_client,
	)
	if err != nil {
		return tradeOrderCreation{}, err
	}
	agg_resp, err := coinprovider.GetTradeTransaction(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		pvk.PublicKey().String(),
		mint,
		utils.ToLamports(sol_amount),
	)
	if err != nil {
		return tradeOrderCreation{}, err
	}
	if agg_resp.ErrorCode != 0 {
		return tradeOrderCreation{}, fmt.Errorf(
			"Received the following error for token %s while generating transaction: %s.",
			mint,
			agg_resp.ErrorMessage,
		)
	}

	totalFees, err := agg_resp.GetTotalFeeAmount()
	if err != nil {
		return tradeOrderCreation{}, err
	}

	solanaAmountAsInt, err := strconv.Atoi(agg_resp.InAmount)
	if err != nil {
		return tradeOrderCreation{}, err
	}

	expectedTokenAmountAsInt, err := strconv.Atoi(agg_resp.OutAmount)
	if err != nil {
		return tradeOrderCreation{}, err
	}

	return tradeOrderCreation{
		Trade: entities.Trade{
			Mint:                mint,
			SolanaAmount:        utils.FromLamports(solanaAmountAsInt),
			TotalFees:           utils.FromLamports(totalFees),
			ExpectedTokenAmount: utils.FromLamports(expectedTokenAmountAsInt),
		},
		Transaction: agg_resp.Transaction,
	}, nil
}

func logTradeSimulation(simulationResult coinprovider.TransactionSimulation) {
	log.Println("Err:", simulationResult.Err)
	for _, logMsg := range simulationResult.Logs {
		log.Println("Log Message:", logMsg)
	}
}
