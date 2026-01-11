-- ============================
-- TRANSACTIONS TABLE
-- ============================

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(64) PRIMARY KEY,
    payment_type VARCHAR(20) NOT NULL,
    payment_mode VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ============================
-- LEDGER ENTRIES TABLE
-- ============================

CREATE TABLE IF NOT EXISTS ledger_entries (
    id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    entry_type VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL,
    source VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transaction
        FOREIGN KEY(transaction_id)
        REFERENCES transactions(id)
);

-- ============================
-- INDEXES
-- ============================

CREATE INDEX IF NOT EXISTS idx_ledger_transaction_id
ON ledger_entries(transaction_id);

CREATE INDEX IF NOT EXISTS idx_transaction_status
ON transactions(status);
