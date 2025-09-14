CREATE TABLE IF NOT EXISTS top_holder (
    mint VARCHAR(44),
    top_5_wallets DOUBLE,
    top_10_wallets DOUBLE,
    top_20_wallets DOUBLE,
    has_single_largest_holder BOOLEAN,
    tracked_at DATETIME,
    CONSTRAINT top_holder_pk PRIMARY KEY (
	mint,
	tracked_at
    )
);
