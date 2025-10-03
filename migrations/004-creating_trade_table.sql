CREATE TABLE IF NOT EXISTS trade (
    mint VARCHAR(44),                        -- Order creation
    operation VARCHAR,                       -- Order creation
    slippage_bps INTEGER,                    -- Order creation
    input_amount_lamports INTEGER,           -- Order creation
    expected_output_amount_lamports INTEGER, -- Order creation
    input_usd_price DOUBLE,                  -- Order creation
    total_fee_lamports INTEGER,              -- disconted from the trade on buy, from wallet on sell.
    expected_token_usd_price DOUBLE,         -- get from database during order creation.
    issued_order_at DATETIME,                -- Order execution
    received_order_response_at DATETIME,     -- Order execution
    executed_output_amount_lamports INTEGER, -- Order execution
    executed_token_usd_price DOUBLE,         -- Order execution
    -- for buy
    -- (sol) input usd price after fees / executed token amount
    -- for sell
    -- (token) input usd price / executed token amount
    CONSTRAINT trade_pk PRIMARY KEY (
	mint,
	operation
    )
);
