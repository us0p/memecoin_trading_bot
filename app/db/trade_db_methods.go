package db

import "context"

func (d *DB) GetOngoingTradesBalance(ctx context.Context) (float64, error) {
	row := d.db.QueryRowContext(
		ctx,
		`SELECT 
			SUM(solana_amount) 
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
