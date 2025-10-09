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

func newTradeOrderCreation(
	order entities.Order,
	slippageBPS,
	inputAmountLamports,
	totalFeeLamport,
	expectedOutputAmountLamports int,
	inputUSDPrice float64,
	transaction string,
) tradeOrderCreation {
	return tradeOrderCreation{
		Trade: entities.Trade{
			Mint:                         order.Mint,
			Operation:                    order.Op,
			InputAmountLamports:          inputAmountLamports,
			TotalFeeLamports:             totalFeeLamport,
			ExpectedOutputAmountLamports: expectedOutputAmountLamports,
			InputUSDPrice:                inputUSDPrice,
		},
		Transaction: transaction,
	}
}

func ExecuteTrade(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
	order entities.Order,
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
		order,
	)
	if err != nil {
		nf_state.RecordError(
			order.Mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}

	if err = signTransaction(pvk, &tradeOrderCreation); err != nil {
		nf_state.RecordError(
			order.Mint,
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
			order.Mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}

	if err = db_client.InsertTrade(ctx, trade); err != nil {
		nf_state.RecordError(
			order.Mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}
}

// Using the signed transaction. Executes the order with Jupiter API.
func executeOrder(
	ctx context.Context,
	http_client *http.Client,
	db_client *db.DB,
	tradeOrder *tradeOrderCreation,
) (entities.Trade, error) {
	tradeOrder.Trade.IssuedOrderAt = time.Now()

	simu, err := coinprovider.SimulateTransactionExecution(
		http_client,
		constants.HELIUS_API_URL,
		tradeOrder.Transaction,
	)
	if err != nil {
		return entities.Trade{}, err
	}

	tradeOrder.Trade.ReceivedOrderResponseAt = time.Now()

	last_price, err := db_client.GetLastPriceForToken(ctx, tradeOrder.Trade.Mint)
	if err != nil {
		return entities.Trade{}, err
	}

	tradeOrder.Trade.ExpectedTokenUSDPrice = last_price

	mk_data, err := coinprovider.GetMarketDataForAddresses(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		[]string{tradeOrder.Trade.Mint},
	)
	if err != nil {
		return entities.Trade{}, err
	}

	// apply math here.
	tradeOrder.Trade.ExecutedTokenUSDPrice = mk_data[0].PriceUsd

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

type orderCb func(*http.Client, *db.DB, solana.PrivateKey, string) (tradeOrderCreation, error)

type orderStrategies map[entities.Operation]orderCb

func newOrderStrategies(buyCb, sellCb orderCb) orderStrategies {
	strategies := make(map[entities.Operation]orderCb)

	strategies[entities.BUY] = buyCb
	strategies[entities.SELL] = sellCb

	return strategies
}

// Create order by calling Jupiter 'order' API.
// Returns a partial TradeOrder with the base64 encoded Transaction.
func getOrder(
	http_client *http.Client,
	db_client *db.DB,
	pvk solana.PrivateKey,
	order entities.Order,
) (tradeOrderCreation, error) {
	orderStrats := newOrderStrategies(buyStrategy, sellStrategy)
	return orderStrats[order.Op](http_client, db_client, pvk, order.Mint)
}

func buyStrategy(
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
		constants.SOL_MINT_ADDRS,
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

	return newTradeOrderCreation(
		entities.Order{Mint: mint, Op: entities.BUY},
		agg_resp.SlippageBps,
		solanaAmountAsInt,
		totalFees,
		expectedTokenAmountAsInt,
		agg_resp.Transaction,
	), nil
}

func sellStrategy(
	http_client *http.Client,
	db_client *db.DB,
	pvk solana.PrivateKey,
	mint string,
) (tradeOrderCreation, error) {
	wallet_holdings, err := coinprovider.GetOnChainWalletHoldings(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		pvk.PublicKey().String(),
	)
	if err != nil {
		return tradeOrderCreation{}, err
	}

	token_holdings, ok := wallet_holdings.Tokens[mint]
	if !ok {
		return tradeOrderCreation{}, fmt.Errorf(
			"Token mint `%s` is not present in wallet holdings",
			mint,
		)
	}
	token_amount, err := strconv.Atoi(token_holdings.Amount)
	if err != nil {
		return tradeOrderCreation{}, err
	}
	agg_resp, err := coinprovider.GetTradeTransaction(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		pvk.PublicKey().String(),
		mint,
		constants.SOL_MINT_ADDRS,
		token_amount,
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

	tokenAmountAsInt, err := strconv.Atoi(agg_resp.InAmount)
	if err != nil {
		return tradeOrderCreation{}, err
	}

	expectedSolanaAmountAsInt, err := strconv.Atoi(agg_resp.OutAmount)
	if err != nil {
		return tradeOrderCreation{}, err
	}

	return newTradeOrderCreation(
		mint,
		tokenAmountAsInt,
		totalFees,
		expectedSolanaAmountAsInt,
		agg_resp.Transaction,
	), nil
}

func logTradeSimulation(simulationResult coinprovider.TransactionSimulation) {
	log.Println("Err:", simulationResult.Err)
	for _, logMsg := range simulationResult.Logs {
		log.Println("Log Message:", logMsg)
	}
}
