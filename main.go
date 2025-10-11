package main

import (
	"log"
	"memecoin_trading_bot/app/db"
	"memecoin_trading_bot/app/entities"
	"memecoin_trading_bot/app/job_scheduler"
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
	job_scheduler := jobscheduler.NewJobSchedulerMap()
	order_chan := make(chan entities.Order)

	job_scheduler.RegisterJob(
		workflows.PullTokens,
		"Pull memecoin tokens",
		jobscheduler.ONE_MINUTE_INTERVAL,
	)
	job_scheduler.RegisterJob(
		workflows.GetTradeOpportunityMarketData,
		"Pull market data",
		jobscheduler.FIVE_SECOND_INTERVAL,
	)
	job_scheduler.RegisterJob(
		workflows.GetTradeOpportunityLargestHolders,
		"Pull largest holders",
		jobscheduler.FIVE_MINUTE_INTERVAL,
	)

	log.Printf("Starting job executor...\n")
	job_scheduler.StartJobExecutor(
		client,
		&db,
		&nf_state,
		order_chan,
	)

	workflows.TradeChannelProcesser(
		client,
		&db,
		&nf_state,
		order_chan,
	)
}
