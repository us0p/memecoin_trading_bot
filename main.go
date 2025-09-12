package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"memecoin_trading_bot/app/coin_provider"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load envs %s", err)
	}

	client := http.DefaultClient
	calls, err := coinprovider.GetGambleTokens(client)
	if err != nil {
		log.Fatal(err)
	}

	for _, call := range calls {
		fmt.Println("mint: ", call.Mint)
		fmt.Println("symbol: ", call.Symbol)
		fmt.Println("created_at: ", call.CreatedAt)
		fmt.Println("-----")
	}

	authorities, err := coinprovider.GetTokenAuthorities(
		client,
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mint authority:", authorities.MintAuthority)
	fmt.Println("Freeze authority:", authorities.FreezeAuthority)
}
