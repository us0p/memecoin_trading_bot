package db

import (
	"context"
	"memecoin_trading_bot/app/constants"
	"memecoin_trading_bot/app/entities"
)

func (d *DB) InsertTrade(ctx context.Context, trade entities.Trade) error {
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO trade
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		trade.Mint,
		trade.IssuedTradeStartAt.Format(constants.JAVASCRIPT_TIME_REP),
		trade.TradeStartedAt.Format(constants.JAVASCRIPT_TIME_REP),
		trade.IssuedTradeEndAt.Format(constants.JAVASCRIPT_TIME_REP),
		trade.TradeEndedAt.Format(constants.JAVASCRIPT_TIME_REP),
		trade.IssuedTradeStartTokenUsdPrice,
		trade.IssuedTradeEndTokenUsdPrice,
		trade.EntryTokenUsdPrice,
		trade.ExitTokenUsdPrice,
		trade.SolanaAmount,
		trade.TotalFees,
		trade.ExpectedTokenAmount,
		trade.ExecutedTokenAmount,
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
		 ORDER BY priced_at DESC
		 LIMIT 1;`,
	)

	var last_price float64
	if err := row.Scan(&last_price); err != nil {
		return last_price, err
	}

	return last_price, nil
}

func (d *DB) GetOngoingTradesBalance(ctx context.Context) (float64, error) {
	row := d.db.QueryRowContext(
		ctx,
		`SELECT 
			COALESCE(SUM(solana_amount), 0)
		 FROM trade 
		 WHERE trade_ended_at IS NULL;`,
	)

	var total_amount float64
	if err := row.Scan(&total_amount); err != nil {
		return total_amount, err
	}
	return total_amount, nil
}

func (d *DB) GetNewTradeMints(ctx context.Context) ([]string, error) {
	rows, err := d.db.QueryContext(
		ctx,
		`SELECT
			mint
		 FROM token
		 LEFT JOIN trade USING(mint)
		 WHERE token.trade_opp IS TRUE
		 	AND trade.mint IS NULL;`,
	)

	mints := make([]string, 0)

	if err != nil {
		return mints, err
	}

	for rows.Next() {
		var mint string
		if err = rows.Scan(&mint); err != nil {
			return mints, err
		}
		mints = append(mints, mint)
	}

	return mints, nil
}
