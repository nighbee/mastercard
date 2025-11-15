-- Script to load transaction data from CSV
-- This will be executed after the table is created
-- Note: CSV file must be accessible from the container

-- First, create a temporary table to match CSV structure exactly
CREATE TEMP TABLE temp_transactions (
    card_no TEXT,
    date TEXT,
    process_date TEXT,
    trx_amount_usd TEXT,
    trx_amount_eur TEXT,
    trx_amount_local TEXT,
    trx_cnt_usd TEXT,
    trx_cnt_eur TEXT,
    trx_cnt_local TEXT,
    interchange_fee TEXT,
    merch_name TEXT,
    agg_merch_name TEXT,
    issuer_code TEXT,
    issuer_country TEXT,
    bin6_code TEXT,
    acquirer_code TEXT,
    acquirer_country TEXT,
    trx_type TEXT,
    trx_direction TEXT,
    mcc TEXT,
    mcc_group TEXT,
    input_mode TEXT,
    wallet_type TEXT,
    product_type TEXT,
    authorization_status TEXT,
    authorization_response_code TEXT,
    location_id TEXT,
    location_city TEXT
);

-- Copy data from CSV file
-- Note: The CSV file path should be relative to the container or mounted volume
\copy temp_transactions FROM '/datasets/transactions.csv' WITH (FORMAT csv, HEADER true, DELIMITER ',');

-- Insert into main transactions table with proper type conversions
INSERT INTO transactions (
    card_no, date, process_date, trx_amount_usd, trx_amount_eur, trx_amount_local,
    trx_cnt_usd, trx_cnt_eur, trx_cnt_local, interchange_fee, merch_name, agg_merch_name,
    issuer_code, issuer_country, bin6_code, acquirer_code, acquirer_country, trx_type,
    trx_direction, mcc, mcc_group, input_mode, wallet_type, product_type,
    authorization_status, authorization_response_code, location_id, location_city
)
SELECT 
    card_no,
    TO_DATE(date, 'YYYY-MM-DD'),
    TO_DATE(process_date, 'YYYY-MM-DD'),
    NULLIF(trx_amount_usd, '')::DECIMAL(15, 2),
    NULLIF(trx_amount_eur, '')::DECIMAL(15, 2),
    NULLIF(trx_amount_local, '')::DECIMAL(15, 2),
    COALESCE(NULLIF(trx_cnt_usd, '')::INTEGER, 0),
    COALESCE(NULLIF(trx_cnt_eur, '')::INTEGER, 0),
    COALESCE(NULLIF(trx_cnt_local, '')::INTEGER, 0),
    NULLIF(interchange_fee, '')::DECIMAL(15, 2),
    merch_name, agg_merch_name,
    issuer_code, issuer_country, bin6_code, acquirer_code, acquirer_country, trx_type,
    trx_direction, mcc, mcc_group, input_mode, wallet_type, product_type,
    authorization_status, authorization_response_code, location_id, location_city
FROM temp_transactions;

-- Drop temporary table
DROP TABLE temp_transactions;

-- Show count of loaded records
SELECT COUNT(*) as total_transactions_loaded FROM transactions;
