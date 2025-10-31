package workflows

import (
	"context"
	"encoding/base64"
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
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
)

type tradeOrderCreation struct {
	Trade       entities.Trade
	Transaction string
	RequestId   string
}

func newTradeOrderCreation(
	order entities.Order,
	slippageBPS,
	inputAmountLamports,
	totalFeeLamport,
	expectedOutputAmountLamports int,
	inputUSDPrice float64,
	transaction,
	requestId string,
) tradeOrderCreation {
	return tradeOrderCreation{
		Trade: entities.Trade{
			Mint:                         order.Mint,
			Operation:                    order.Op,
			SlippageBPS:                  slippageBPS,
			InputAmountLamports:          inputAmountLamports,
			TotalFeeLamports:             totalFeeLamport,
			ExpectedOutputAmountLamports: expectedOutputAmountLamports,
			InputUSDPrice:                inputUSDPrice,
		},
		Transaction: transaction,
		RequestId:   requestId,
	}
}

type OrderStatus string

const (
	processing OrderStatus = "processing"
	received               = "received"
)

type TransactionProcessing struct {
	OrderChan          chan entities.Order
	OrderProcessingMap map[entities.Order]OrderStatus
	mut                sync.RWMutex
}

func NewTransactionProcessing() TransactionProcessing {
	return TransactionProcessing{
		make(chan entities.Order),
		make(map[entities.Order]OrderStatus),
		sync.RWMutex{},
	}
}

func (t *TransactionProcessing) IssueOrder(order entities.Order) {
	if _, ok := (*t).OrderProcessingMap[order]; ok {
		return
	}

	t.mut.Lock()
	defer t.mut.Unlock()
	(*t).OrderProcessingMap[order] = received
	t.OrderChan <- order
}

func (t *TransactionProcessing) FulfillOrder(order entities.Order) {
	t.mut.Lock()
	defer t.mut.Unlock()
	delete((*t).OrderProcessingMap, order)
}

func TradeChannelProcesser(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
	tp *TransactionProcessing,
) {
	for order := range tp.OrderChan {
		if tp.OrderProcessingMap[order] == processing {
			return
		}
		tp.OrderProcessingMap[order] = processing
		log.Printf("Received new trade opportunity for address: %s, %s\n", order.Mint, order.Op)
		go executeTrade(
			http_client,
			db_client,
			nf_state,
			order,
			tp,
		)
	}
}

func executeTrade(
	http_client *http.Client,
	db_client *db.DB,
	nf_state *notification.Notifications,
	order entities.Order,
	tp *TransactionProcessing,
) {
	log.Println("Starting trade execution")
	ctx := context.Background()

	log.Println("Getting Prv")
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

	log.Println("Getting order...")
	tradeOrderCreation, wallet_ballance, err := getOrder(
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

	log.Println("Signing transaction...")
	if err = signTransaction(pvk, &tradeOrderCreation); err != nil {
		nf_state.RecordError(
			order.Mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}

	log.Println("Executing order")
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

	log.Println("Inserting trade...")
	if err = db_client.InsertTrade(ctx, trade); err != nil {
		nf_state.RecordError(
			order.Mint,
			notification.ExecuteTrade,
			err,
			notification.Core,
		)
		return
	}

	tp.FulfillOrder(order)

	nf_state.RecordTradeExecution(
		db_client,
		order.Mint,
		order.Op,
		wallet_ballance,
	)
}

// Using the signed transaction. Executes the order with Jupiter API.
func executeOrder(
	ctx context.Context,
	http_client *http.Client,
	db_client *db.DB,
	tradeOrder *tradeOrderCreation,
) (entities.Trade, error) {
	tradeOrder.Trade.IssuedOrderAt = time.Now()

	order, err := coinprovider.ExecuteTransaction(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		tradeOrder.Transaction,
		tradeOrder.RequestId,
	)
	if err != nil {
		return entities.Trade{}, err
	}
	exec_output_amout_lamp, err := strconv.Atoi(order.OutputAmountResut)
	if err != nil {
		return entities.Trade{}, err
	}
	tradeOrder.Trade.ExecutedOutputAmountLamports = exec_output_amout_lamp

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

	tradeOrder.Trade.ExecutedTokenUSDPrice = mk_data[0].PriceUsd

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

	signedTxBytes, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	*&tradeOrder.Transaction = base64.StdEncoding.EncodeToString(signedTxBytes)
	return nil
}

type orderCb func(*http.Client, *db.DB, solana.PrivateKey, string) (tradeOrderCreation, float64, error)

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
) (tradeOrderCreation, float64, error) {
	orderStrats := newOrderStrategies(buyStrategy, sellStrategy)
	return orderStrats[order.Op](http_client, db_client, pvk, order.Mint)
}

func buyStrategy(
	http_client *http.Client,
	db_client *db.DB,
	pvk solana.PrivateKey,
	mint string,
) (tradeOrderCreation, float64, error) {
	sol_amount, wallet_balance, err := riskmanagement.GetTradeAmount(
		http_client,
		db_client,
	)
	if err != nil {
		return tradeOrderCreation{}, wallet_balance, err
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
		return tradeOrderCreation{}, wallet_balance, err
	}
	if agg_resp.ErrorCode != 0 {
		return tradeOrderCreation{}, wallet_balance, fmt.Errorf(
			"Received the following error for token %s while generating transaction: %s.",
			mint,
			agg_resp.ErrorMessage,
		)
	}

	totalFees, err := agg_resp.GetTotalFeeAmount()
	if err != nil {
		return tradeOrderCreation{}, wallet_balance, err
	}

	solanaAmountAsInt, err := strconv.Atoi(agg_resp.InAmount)
	if err != nil {
		return tradeOrderCreation{}, wallet_balance, err
	}

	expectedTokenAmountAsInt, err := strconv.Atoi(agg_resp.OutAmount)
	if err != nil {
		return tradeOrderCreation{}, wallet_balance, err
	}

	return newTradeOrderCreation(
		entities.Order{Mint: mint, Op: entities.BUY},
		agg_resp.SlippageBps,
		solanaAmountAsInt,
		totalFees,
		expectedTokenAmountAsInt,
		agg_resp.InUsdValue,
		agg_resp.Transaction,
		agg_resp.RequestId,
	), wallet_balance, nil
}

func sellStrategy(
	http_client *http.Client,
	db_client *db.DB,
	pvk solana.PrivateKey,
	mint string,
) (tradeOrderCreation, float64, error) {
	wallet_holdings, err := coinprovider.GetOnChainWalletHoldings(
		http_client,
		constants.JUPITER_ULTRA_API_URL,
		pvk.PublicKey().String(),
	)
	log.Println(constants.JUPITER_ULTRA_API_URL, pvk.PublicKey().String())
	if err != nil {
		return tradeOrderCreation{}, 0, err
	}

	token_holdings, ok := wallet_holdings.Tokens[mint]
	if !ok {
		return tradeOrderCreation{}, 0, fmt.Errorf(
			"Token mint `%s` is not present in wallet holdings",
			mint,
		)
	}
	token_amount, err := strconv.Atoi(token_holdings[0].Amount)
	if err != nil {
		return tradeOrderCreation{}, 0, err
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
		return tradeOrderCreation{}, 0, err
	}
	if agg_resp.ErrorCode != 0 {
		return tradeOrderCreation{}, 0, fmt.Errorf(
			"Received the following error for token %s while generating transaction: %s.",
			mint,
			agg_resp.ErrorMessage,
		)
	}

	totalFees, err := agg_resp.GetTotalFeeAmount()
	if err != nil {
		return tradeOrderCreation{}, 0, err
	}

	tokenAmountAsInt, err := strconv.Atoi(agg_resp.InAmount)
	if err != nil {
		return tradeOrderCreation{}, 0, err
	}

	expectedSolanaAmountAsInt, err := strconv.Atoi(agg_resp.OutAmount)
	if err != nil {
		return tradeOrderCreation{}, 0, err
	}

	return newTradeOrderCreation(
		entities.Order{
			Mint: mint,
			Op:   entities.SELL,
		},
		agg_resp.SlippageBps,
		tokenAmountAsInt,
		totalFees,
		expectedSolanaAmountAsInt,
		agg_resp.InUsdValue,
		agg_resp.Transaction,
		agg_resp.RequestId,
	), 0, nil
}

func logTradeSimulation(simulationResult coinprovider.TransactionSimulation) {
	log.Println("Err:", simulationResult.Err)
	for _, logMsg := range simulationResult.Logs {
		log.Println("Log Message:", logMsg)
	}
}
