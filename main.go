package main

import (
	"log"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/notification"

	"memecoin_trading_bot/app/workflows"

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
	nf_state := notification.NewNotificationState()

	workflows.PullTokens(
		client,
		&db,
		&nf_state,
	)

	workflows.GetTradeOpportunityMarketData(
		client,
		&db,
		&nf_state,
	)

	workflows.GetTradeOpportunityLargestHolders(
		client,
		&db,
		&nf_state,
	)

	workflows.ExecuteTrade(
		client,
		&db,
		&nf_state,
		"5zCETicUCJqJ5Z3wbfFPZqtSpHPYqnggs1wX7ZRpump",
	)

	nf_state.SendNotifications(
		client,
		constants.TELEGRAM_API_URL,
	)
}
