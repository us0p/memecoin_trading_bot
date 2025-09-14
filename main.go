package main

import (
	"log"
	"memecoin_trading_bot/app/db"
	//"net/http"

	"github.com/joho/godotenv"
)

const memescan_api_url = "https://memescan.app/api/calls?sort=recent&limit=10&offset=0&type=gamble"
const jupiter_ultra_api_url = "https://lite-api.jup.ag/ultra/v1"
const helius_api_url = "https://mainnet.helius-rpc.com"

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

	//client := http.DefaultClient
}
