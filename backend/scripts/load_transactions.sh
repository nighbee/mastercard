#!/bin/bash
# Script to load transaction data from CSV into PostgreSQL
# Usage: ./load_transactions.sh

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values if not in .env
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-mastercard_user}
DB_NAME=${DB_NAME:-mastercard_db}
PGPASSWORD=${DB_PASSWORD:-mastercard_pass}
CSV_FILE=${CSV_FILE:-../../doc/datasets/transactions.csv}

echo "Loading transaction data from CSV..."
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo "CSV File: $CSV_FILE"
echo ""

# Check if CSV file exists
if [ ! -f "$CSV_FILE" ]; then
    echo "Error: CSV file not found at $CSV_FILE"
    exit 1
fi

# Create temporary table
echo "Creating temporary table..."
PGPASSWORD=$PGPASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
CREATE TEMP TABLE temp_transactions (
    transaction_id TEXT,
    transaction_timestamp TEXT,
    card_id TEXT,
    expiry_date TEXT,
    issuer_bank_name TEXT,
    merchant_id TEXT,
    merchant_mcc TEXT,
    mcc_category TEXT,
    merchant_city TEXT,
    transaction_type TEXT,
    transaction_amount_kzt TEXT,
    original_amount TEXT,
    transaction_currency TEXT,
    acquirer_country_iso TEXT,
    pos_entry_mode TEXT,
    wallet_type TEXT,
    index_level_0 TEXT
);
EOF

# Copy data from CSV
echo "Copying data from CSV..."
PGPASSWORD=$PGPASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
\copy temp_transactions FROM '$CSV_FILE' WITH (FORMAT csv, HEADER true, DELIMITER ',');
EOF

# Insert into main table
echo "Inserting data into transactions table..."
PGPASSWORD=$PGPASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
INSERT INTO transactions (
    transaction_id,
    transaction_timestamp,
    card_id,
    expiry_date,
    issuer_bank_name,
    merchant_id,
    merchant_mcc,
    mcc_category,
    merchant_city,
    transaction_type,
    transaction_amount_kzt,
    original_amount,
    transaction_currency,
    acquirer_country_iso,
    pos_entry_mode,
    wallet_type,
    index_level_0
)
SELECT 
    transaction_id::UUID,
    transaction_timestamp::TIMESTAMP,
    NULLIF(card_id, '')::INTEGER,
    NULLIF(expiry_date, ''),
    NULLIF(issuer_bank_name, ''),
    NULLIF(merchant_id, '')::INTEGER,
    NULLIF(merchant_mcc, '')::INTEGER,
    NULLIF(mcc_category, ''),
    NULLIF(merchant_city, ''),
    NULLIF(transaction_type, ''),
    NULLIF(transaction_amount_kzt, '')::DECIMAL(15, 2),
    NULLIF(original_amount, '')::DECIMAL(15, 2),
    NULLIF(transaction_currency, ''),
    NULLIF(acquirer_country_iso, ''),
    NULLIF(pos_entry_mode, ''),
    NULLIF(wallet_type, ''),
    NULLIF(index_level_0, '')::INTEGER
FROM temp_transactions
WHERE transaction_id IS NOT NULL AND transaction_id != '';

DROP TABLE temp_transactions;

SELECT COUNT(*) as total_transactions_loaded FROM transactions;
EOF

echo ""
echo "✓ Transaction data loaded successfully!"

