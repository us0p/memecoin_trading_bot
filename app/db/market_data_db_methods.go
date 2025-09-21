package db

import (
	"context"
	"fmt"
	coinprovider "memecoin_trading_bot/app/coin_provider"
	"memecoin_trading_bot/app/entities"
	"strings"
	"time"
)

func (d *DB) GetLatestSupplyForToken(ctx context.Context, mint string) (float64, error) {
	row := d.db.QueryRow(`
	SELECT
		total_supply
	FROM market_data
	ORDER BY priced_at DESC
	LIMIT 1;
	`)

	var total_supply float64
	if err := row.Scan(&total_supply); err != nil {
		return total_supply, err
	}

	return total_supply, nil
}

func (d *DB) InsertMarketDataBulk(
	ctx context.Context,
	marketData []coinprovider.MarketData,
) error {
	placeholders := make([]string, len(marketData))
	args := make([]any, 0, len(marketData)*16)
	for idx, mkd := range marketData {
		placeholders[idx] = "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		args = append(args,
			mkd.Mint,
			mkd.TotalSupply,
			time.Now(),
			mkd.PriceUsd,
			mkd.Liquidity,
			mkd.HolderCount,
			mkd.MarketCap,
			mkd.Stats1h.VolumeChange,
			mkd.Stats1h.BuyVolume,
			mkd.Stats1h.SellVolume,
			mkd.Stats1h.BuyOrganicVolume,
			mkd.Stats1h.SellOrganicVolume,
			mkd.Stats1h.NumBuys,
			mkd.Stats1h.NumSells,
			mkd.Stats1h.NumTraders,
			mkd.Stats1h.NumOrganicBuyers,
			mkd.Stats1h.NumNetBuyers,
		)
	}
	query := fmt.Sprintf(`
		INSERT INTO market_data
		VALUES %s;
	`, strings.Join(placeholders, ","))

	_, err := d.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (d *DB) InsertTopHolderBulk(
	ctx context.Context,
	top_holders []entities.TopHolder,
) error {
	placeholders := make([]string, len(top_holders))
	args := make([]any, 0, len(top_holders)*6)
	for idx, th := range top_holders {
		placeholders[idx] = "(?,?,?,?,?,?)"
		args = append(args,
			th.Mint,
			th.Top5Wallets,
			th.Top10Wallets,
			th.Top20Wallets,
			th.HasSingleLargestHolder,
			time.Now(),
		)
	}
	query := fmt.Sprintf(`
		INSERT INTO top_holder
		VALUES %s;
	`, strings.Join(placeholders, ","))

	_, err := d.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}
