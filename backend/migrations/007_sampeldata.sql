-- This script inserts 100 random sample records into the 'transactions' table.
-- It uses generate_series to loop 100 times and random functions to create varied data.

INSERT INTO transactions (
    card_no,
    date,
    process_date,
    trx_amount_usd,
    merch_name,
    issuer_country,
    trx_type,
    mcc,
    mcc_group,
    location_city,
    authorization_status
)
SELECT
    -- Generate a random card number
    '5'||LPAD((100000000000000 + random() * 899999999999999)::BIGINT::TEXT, 15, '0'),

    -- Generate a random timestamp within the last year
    (NOW() - (random() * 365 * '1 day'::interval))::DATE,
    (NOW() - (random() * 364 * '1 day'::interval))::DATE,

    -- Generate a random transaction amount in USD
    (random() * 2000 + 10)::DECIMAL(15, 2),
    
    -- Pick a random merchant name
    (ARRAY['Amazon', 'Starbucks', 'Netflix', 'Shell', 'Walmart', 'Apple Store'])[ (1 + random() * 5)::INTEGER ],

    -- Generate a random 4-digit MCC code
    (ARRAY['USA', 'GBR', 'KAZ', 'DEU', 'FRA'])[ (1 + random() * 4)::INTEGER ],

    -- Pick a random MCC category
    (ARRAY['POS', 'ONLINE', 'ATM', 'TRANSFER'])[ (1 + random() * 3)::INTEGER ],

    -- Pick a random city
    (5000 + random() * 4999)::INTEGER::TEXT,

    -- Pick a random transaction type
    (ARRAY['Retail', 'Groceries', 'Restaurants', 'Travel', 'Entertainment'])[ (1 + random() * 4)::INTEGER ],

    -- Pick a random city
    (ARRAY['Almaty', 'Astana', 'Shymkent', 'New York', 'London'])[ (1 + random() * 4)::INTEGER ],

    -- Pick a random status
    (ARRAY['approved', 'declined'])[ (1 + random())::INTEGER ]

FROM generate_series(1, 100) AS s(id); -- This generates the 100 rows
