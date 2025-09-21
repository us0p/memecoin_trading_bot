package main

import (
	"fmt"
	"log"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	//"memecoin_trading_bot/app/notification"
	//"memecoin_trading_bot/app/workflows"

	"net/http"

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

	dt, err := coinprovider.GetTradeTransaction(
		client,
		constants.JUPITER_ULTRA_API_URL,
		"5zCETicUCJqJ5Z3wbfFPZqtSpHPYqnggs1wX7ZRpump",
		1000000000,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", dt)
}
