CREATE TABLE IF NOT EXISTS token (
    mint VARCHAR(44) PRIMARY KEY,
    symbol TEXT NOT NULL,
    mint_enabled BOOLEAN NOT NULL,
    freeze_enabled BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL,
    trade_opp BOOLEAN,
    twitter TEXT,
    site TEXT,
    telegram TEXT
);
