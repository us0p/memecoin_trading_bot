package db

import (
	"context"
	"memecoin_trading_bot/app/entities"
	"strings"
)

func (d *DB) QueryActiveTokens(ctx context.Context) ([]entities.Token, error) {
	var tokens []entities.Token

	rows, err := d.db.QueryContext(ctx, `SELECT * FROM token WHERE trade_opp IS TRUE;`, nil)
	if err != nil {
		return tokens, err
	}

	for rows.Next() {
		var token entities.Token
		if err = rows.Scan(
			&token.Mint,
			&token.Symbol,
			&token.MintEnabled,
			&token.FreezeEnabled,
			&token.CreatedAt,
			&token.TradeOpp,
			&token.Twitter,
			&token.Site,
			&token.Telegram,
		); err != nil {
			return tokens, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (d *DB) QueryExistingTokensFromSlice(ctx context.Context, mints []string) ([]string, error) {
	var newMints []string

	rows, err := d.db.QueryContext(ctx, `
		SELECT
			mint
		FROM token
		WHERE mint IN (?);
	`, strings.Join(mints, ","))
	if err != nil {
		return newMints, err
	}

	for rows.Next() {
		var mint string
		if err = rows.Scan(&mint); err != nil {
			return newMints, err
		}
		newMints = append(newMints, mint)
	}

	return newMints, nil
}

func (d *DB) InsertToken(ctx context.Context, token entities.Token) (entities.Token, error) {
	var newToken entities.Token
	query, err := d.db.Prepare(`
		INSERT INTO token
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING *;
	`)
	if err != nil {
		return newToken, err
	}

	row := query.QueryRowContext(
		ctx,
		token.Mint,
		token.Symbol,
		token.MintEnabled,
		token.FreezeEnabled,
		token.CreatedAt,
		token.TradeOpp,
		token.Twitter,
		token.Site,
		token.Telegram,
	)

	err = row.Scan(
		&newToken.Mint,
		&newToken.Symbol,
		&newToken.MintEnabled,
		&newToken.FreezeEnabled,
		&newToken.CreatedAt,
		&newToken.TradeOpp,
		&newToken.Twitter,
		&newToken.Site,
		&newToken.Telegram,
	)
	if err != nil {
		return newToken, err
	}

	return newToken, nil
}
