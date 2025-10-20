package db

import (
	"context"
	"database/sql"
	"errors"
	"memecoin_trading_bot/app/entities"
)

func (d *DB) FulfillTransaction(
	ctx context.Context,
	order entities.Order,
) error {
	_, err := d.db.ExecContext(
		ctx,
		`DELETE FROM trade_transaction_processing 
		 WHERE ming = ? 
			AND operation = ?;`,
		order.Mint,
		order.Op,
	)
	return err
}

func (d *DB) GetTradeTransactionProcessing(
	ctx context.Context,
	order entities.Order,
) (bool, error) {
	row := d.db.QueryRowContext(
		ctx,
		`SELECT 
			mint 
	 	 FROM trade_transaction_processing 
		 WHERE mint = ? 
		 	AND operation = ?
			AND status <> 'processing';`,
		order.Mint,
		order.Op,
	)

	var mint string
	if err := row.Scan(&mint); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return true, nil
	}
	return true, nil
}

func (d *DB) InsertTradeProcessing(
	ctx context.Context,
	order entities.Order,
) error {
	_, err := d.db.ExecContext(
		ctx,
		`INSERT INTO trade_transaction_processing 
		 VALUES(?, ?, 'processing')`,
		order.Mint,
		order.Op,
	)
	return err
}

func (d *DB) GetTradeNotificationData(
	ctx context.Context,
	mint string,
	op entities.Operation,
) (entities.TradeNotificationData, error) {
	row := d.db.QueryRowContext(
		ctx,
		`
			SELECT	
				symbol,
				received_order_response_at,
				input_usd_price,
				executed_token_usd_price,
				input_amount_lamports
			FROM trade
			JOIN token using(mint)
			WHERE mint = ?
				AND trade.operation = ?;
		`,
		mint,
		op,
	)

	var trade_notification_data entities.TradeNotificationData
	err := row.Scan(
		&trade_notification_data.Symbol,
		&trade_notification_data.ReceivedOrderResponseAt,
		&trade_notification_data.InputUSDPrice,
		&trade_notification_data.ExecutedTokenUSDPrice,
		&trade_notification_data.InputAmountLamports,
	)
	return trade_notification_data, err
}

// pull coin only dispatches buy orders for new tokens that are trade opp.
// market data pull data for the latest 10 trade op tokens, it doesn't matter if they have open trades.
// this function should only return data for tokens that are not closed bc it's retun value
// is used to issue sell orders.
func (d *DB) GetOpenTradeForMint(ctx context.Context, mint string) (entities.Trade, error) {
	row := d.db.QueryRowContext(
		ctx,
		`
			SELECT
				*
			FROM trade
			WHERE mint = ?
			GROUP BY mint
			HAVING COUNT(*) = 1;
		`,
		mint,
	)

	var trade entities.Trade
	err := row.Scan(
		&trade.Mint,
		&trade.Operation,
		&trade.SlippageBPS,
		&trade.InputAmountLamports,
		&trade.ExpectedOutputAmountLamports,
		&trade.InputUSDPrice,
		&trade.TotalFeeLamports,
		&trade.ExpectedTokenUSDPrice,
		&trade.IssuedOrderAt,
		&trade.ReceivedOrderResponseAt,
		&trade.ExecutedOutputAmountLamports,
		&trade.ExecutedTokenUSDPrice,
	)
	if err != nil {
		return trade, err
	}
	return trade, nil
}

func (d *DB) CheckExistingTradeForToken(ctx context.Context, mint string) (bool, error) {
	row := d.db.QueryRowContext(
		ctx,
		`
			SELECT
				mint
			FROM trade
			WHERE mint = ?;
		`,
		mint,
	)

	var existing_token_mint string
	if err := row.Scan(&existing_token_mint); errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}
	return false, nil
}

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
