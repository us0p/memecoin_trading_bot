package db

import (
	"context"
	"memecoin_trading_bot/app/entities"
)

func (d *DB) InsertTrade(ctx context.Context, trade entities.Trade) error {
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO trade
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		trade.Mint,
		trade.Operation,
		trade.SlippageBPS,
		trade.InputAmountLamports,
		trade.ExpectedOutputAmountLamports,
		trade.InputUSDPrice,
		trade.TotalFeeLamports,
		trade.ExpectedTokenUSDPrice,
		trade.IssuedOrderAt,
		trade.ReceivedOrderResponseAt,
		trade.ExecutedOutputAmountLamports,
		trade.ExecutedTokenUSDPrice,
	)

	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetLastPriceForToken(ctx context.Context, mint string) (float64, error) {
	row := d.db.QueryRowContext(
		ctx,
		`SELECT
			price_usd
		 FROM market_data
		 WHERE price_usd IS NOT NULL
		 	AND price_usd != 0.0
			AND mint = ?
		 ORDER BY priced_at DESC
		 LIMIT 1;`,
		mint,
	)

	var last_price float64
	if err := row.Scan(&last_price); err != nil {
		return last_price, err
	}

	return last_price, nil
}

func (d *DB) GetOngoingTradesBalanceLamports(ctx context.Context) (int, error) {
	row := d.db.QueryRowContext(
		ctx,
		`WITH total_buys_aggregate AS (
		 	SELECT 
		 	       mint,
			       SUM(input_amount_lamports) total
		 	FROM trade 
		 	WHERE operation = 'BUY'
		 	GROUP BY mint
		 	HAVING COUNT(*) = 1
		)
		SELECT
			COALESCE(SUM(total), 0)
		FROM total_buys_aggregate;`,
	)

	var total_amount int
	if err := row.Scan(&total_amount); err != nil {
		return total_amount, err
	}
	return total_amount, nil
}
