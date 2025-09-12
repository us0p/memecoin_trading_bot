package main

import (
	//"fmt"
	"log"
	//"memecoin_trading_bot/app/coin_provider"
	//"net/http"

	"github.com/joho/godotenv"
)

const memescan_api_url = "https://memescan.app/api/calls?sort=recent&limit=10&offset=0&type=gamble"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load envs %s", err)
	}

	//client := http.DefaultClient
}
