#	field_name	field_description
1	card_no	Masked or encrypted card number used in the transaction.
2	date	Date when the transaction was initiated by the cardholder.
3	process date	Date when the transaction was processed by the payment system.
4	trx_amount_usd	Transaction amount converted to US Dollars.
5	trx_amount_eur	Transaction amount converted to Euros.
6	trx_amount_local	Transaction amount in the local currency of the transaction.
7	trx_cnt_usd	Number of transactions in USD.
8	trx_cnt_eur	Number of transactions in EUR.
9	trx_cnt_local	Number of transactions in local currency.
10	interchange_fee	Fee charged by the card network for processing the transaction.
11	merch_name	Name of the merchant where the transaction occurred.
12	agg_merch_name	Aggregated or standardized merchant name for grouping similar merchants.
13	issuer_code	Code identifying the bank that issued the card.
14	issuer_country	Country of the issuing bank.
15	bin6_code	First 6 digits of the card number (Bank Identification Number). Useful for identifyign the type of card for campaign analysis.
16	acquirer_code	Code identifying the acquiring bank or payment processor.
17	acquirer_country	Country of the acquiring bank.
18	trx_type	Type of transaction (e.g., POS, ATM, Transfer).
19	trx_direction	plus - outgoing, minus - incoming
20	mcc	Merchant Category Code
21	mcc_group	Grouping of MCCs into broader business categories (e.g., retail, travel).
22	input_mode	Method used to input the card (e.g., chip, swipe, contactless, manual, card on file).
23	wallet_type	Type of digital wallet used (e.g., Apple Pay, Google Pay).
24	product_type	Type of card product (e.g., mass, premium, debit, credit).
25	authorization status	Status of the transaction authorization (e.g., approved, declined).
26	authorization_response_code	Code returned by the issuer indicating the reason for approval or decline.
27	location_id	Internal identifier for the transaction location.
28	location_city	City where the transaction took place.
