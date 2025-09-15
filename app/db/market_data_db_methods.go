package db

import (
	"context"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"time"
)

func (d *DB) InsertMarketData(ctx context.Context, marketData coinprovider.MarketData) error {
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO market_data
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`,
		marketData.Mint,
		time.Now(),
		marketData.PriceUsd,
		marketData.Liquidity,
		marketData.HolderCount,
		marketData.MarketCap,
		marketData.Stats1h.VolumeChange,
		marketData.Stats1h.BuyVolume,
		marketData.Stats1h.SellVolume,
		marketData.Stats1h.BuyOrganicVolume,
		marketData.Stats1h.SellOrganicVolume,
		marketData.Stats1h.NumBuys,
		marketData.Stats1h.NumSells,
		marketData.Stats1h.NumTraders,
		marketData.Stats1h.NumOrganicBuyers,
		marketData.Stats1h.NumNetBuyers,
	)

	if err != nil {
		return err
	}

	return nil
}
