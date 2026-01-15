-- ============================
-- TRANSACTIONS TABLE
-- ============================

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(64) PRIMARY KEY,              -- TM transaction id

    user_ref_id VARCHAR(64) NOT NULL,        -- reference to User Service
    source_ref_id VARCHAR(64) NOT NULL,      -- account/card reference
    destination_ref_id VARCHAR(64),          -- merchant/bank reference

    payment_type VARCHAR(20) NOT NULL,       -- IMMEDIATE, NEFT, CARD
    payment_mode VARCHAR(20) NOT NULL,       -- IMPS, UPI, NEFT, CARD

    status VARCHAR(20) NOT NULL,             -- INITIATED, PENDING, AUTHORIZED, COMPLETED, RELEASED, FAILED
    dc_flag CHAR(1) NOT NULL,                -- D = Debit, C = Credit
    amount BIGINT NOT NULL,
    currency VARCHAR(10) DEFAULT 'INR',

    network_txn_id VARCHAR(64),              -- Visa/Mastercard/UPI/NEFT reference
    gateway_txn_id VARCHAR(64),              -- Razorpay/Stripe etc

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================
-- LEDGER ENTRIES TABLE
-- ============================
CREATE TABLE IF NOT EXISTS ledger_entries (
    id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    account_ref_id VARCHAR(64) NOT NULL,     -- whose balance impacted
    dc_flag CHAR(1) NOT NULL,                -- D = Debit, C = Credit
    entry_type VARCHAR(20) NOT NULL,         -- AUTH, SETTLEMENT, RELEASE, REVERSAL, DEBIT, CREDIT
    amount BIGINT NOT NULL,
    source VARCHAR(20) NOT NULL,             -- API, EVENT, SYSTEM
    msg VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT fk_tx FOREIGN KEY(transaction_id) REFERENCES transactions(id)
);
-- ============================
-- ============================
-- INDEXES
-- ============================

CREATE INDEX IF NOT EXISTS idx_ledger_transaction_id
ON ledger_entries(transaction_id);

CREATE INDEX IF NOT EXISTS idx_transaction_status
ON transactions(status);

CREATE INDEX IF NOT EXISTS idx_ledger_transaction_id
ON ledger_entries(transaction_id);

CREATE INDEX IF NOT EXISTS idx_transaction_status
ON transactions(status);