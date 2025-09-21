CREATE TABLE IF NOT EXISTS market_data (
    mint VARCHAR(44),
    total_supply DOUBLE,
    priced_at DATETIME,
    price_usd DOUBLE,
    liquidity DOUBLE,
    holder_count INT,
    market_cap DOUBLE,
    volume_change_1h DOUBLE,
    buy_volume_1h DOUBLE,
    sell_volume_1h DOUBLE,
    buy_organic_volume_1h DOUBLE,
    sell_organic_volume_1h DOUBLE,
    num_buys_1h INT,
    num_sells_1h INT,
    num_traders_1h INT,
    num_organic_buyers_1h INT,
    num_net_buyers_1h INT,
    CONSTRAINT market_data_pk PRIMARY KEY (
	mint,
	priced_at
    )
);
