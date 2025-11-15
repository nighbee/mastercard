-- Create transactions table matching CSV structure
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    card_no VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    process_date DATE NOT NULL,
    trx_amount_usd DECIMAL(15, 2),
    trx_amount_eur DECIMAL(15, 2),
    trx_amount_local DECIMAL(15, 2),
    trx_cnt_usd INTEGER DEFAULT 0,
    trx_cnt_eur INTEGER DEFAULT 0,
    trx_cnt_local INTEGER DEFAULT 0,
    interchange_fee DECIMAL(15, 2),
    merch_name VARCHAR(255),
    agg_merch_name VARCHAR(255),
    issuer_code VARCHAR(50),
    issuer_country VARCHAR(100),
    bin6_code VARCHAR(6),
    acquirer_code VARCHAR(50),
    acquirer_country VARCHAR(100),
    trx_type VARCHAR(50),
    trx_direction VARCHAR(10) CHECK (trx_direction IN ('plus', 'minus')),
    mcc VARCHAR(10),
    mcc_group VARCHAR(100),
    input_mode VARCHAR(50),
    wallet_type VARCHAR(50),
    product_type VARCHAR(50),
    authorization_status VARCHAR(50),
    authorization_response_code VARCHAR(10),
    location_id VARCHAR(100),
    location_city VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_merch_name ON transactions(merch_name);
CREATE INDEX IF NOT EXISTS idx_transactions_location_city ON transactions(location_city);
CREATE INDEX IF NOT EXISTS idx_transactions_mcc ON transactions(mcc);
