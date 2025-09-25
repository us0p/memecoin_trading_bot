CREATE TABLE IF NOT EXISTS trade (
    mint VARCHAR(44) PRIMARY KEY, -- order creation
    issued_trade_start_at DATETIME, -- order execution
    trade_started_at DATETIME, -- order execution
    issued_trade_end_at DATETIME, -- order execution
    trade_ended_at DATETIME, -- order execution
    issued_trade_start_token_usd_price DOUBLE, -- last mk data entry for token, AFTER order execution
    issued_trade_end_token_usd_price DOUBLE, -- last mk data entry for token, AFTER order execution
    entry_token_usd_price DOUBLE, -- make call for price
    exit_token_usd_price DOUBLE, -- make call for price
    solana_amount DOUBLE, -- order creation
    total_fees DOUBLE, -- order creation
    exepected_token_amount, -- order creation
    executed_token_amount DOUBLE -- order execution
);
