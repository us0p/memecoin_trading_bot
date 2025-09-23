package main

import (
	"encoding/json"
	"fmt"
	"log"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/riskmanagement"
	"memecoin_trading_bot/app/utils"

	//"memecoin_trading_bot/app/workflows"

	"net/http"

	"github.com/gagliardetto/solana-go"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load envs %s", err)
	}

	db, err := db.NewDB("assets.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrate("migrations")
	if err != nil {
		log.Fatal(err)
	}

	client := http.DefaultClient
	//nf_state := notification.NewNotificationState()

	trade_amount, err := riskmanagement.GetTradeAmount(
		client,
		&db,
	)
	if err != nil {
		log.Fatal(err)
	}

	wallet_pvk, err := utils.GetPrvKey()
	if err != nil {
		log.Fatal(err)
	}
	dt, err := coinprovider.GetTradeTransaction(
		client,
		constants.JUPITER_ULTRA_API_URL,
		wallet_pvk.PublicKey().String(),
		"5zCETicUCJqJ5Z3wbfFPZqtSpHPYqnggs1wX7ZRpump",
		utils.ToLamports(trade_amount),
	)

	if err != nil {
		log.Fatal(err)
	}

	tx, err := solana.TransactionFromBase64(dt.Transaction)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet_pvk.PublicKey().Equals(key) {
			return &wallet_pvk
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	signed_tx_base64, err := tx.ToBase64()
	if err != nil {
		log.Fatal(err)
	}

	simulated_tx, err := coinprovider.SimulateTransactionExecution(
		client,
		constants.HELIUS_API_URL,
		signed_tx_base64,
	)
	if err != nil {
		log.Fatal(err)
	}

	byt, _ := json.Marshal(simulated_tx)
	fmt.Println(string(byt))
}
