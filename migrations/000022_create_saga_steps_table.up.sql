CREATE TABLE saga_steps (
    id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    tx_state VARCHAR(20),
    step_name VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_saga_tx
        FOREIGN KEY(transaction_id)
        REFERENCES transactions(id)
);
