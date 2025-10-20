CREATE TABLE IF NOT EXISTS trade_transaction_processing (
    mint VARCHAR(44),
    operation VARCHAR,
    status VARCHAR,
    constraint trade_transaction_processing_pk PRIMARY KEY (
	mint,
	operation
    )
);
