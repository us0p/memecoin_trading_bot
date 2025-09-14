CREATE TABLE IF NOT EXISTS trade (
    mint VARCHAR(44) PRIMARY KEY,
    issued_trade_start_at DATETIME,
    trade_started_at DATETIME,
    issued_trade_end_at DATETIME,
    trade_ended_at DATETIME,
    issued_trade_end_usd_price DOUBLE,
    issued_trade_start_usd_price DOUBLE,
    entry_usd_price DOUBLE,
    exit_usd_price DOUBLE,
    solana_amount DOUBLE,
    executed_solana_amount DOUBLE,
    total_fees DOUBLE
);
