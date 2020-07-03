package schema

// TradeLogsSchema is postgres schema for tradelog
const TradeLogsSchema = `
CREATE TABLE IF NOT EXISTS "users" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL,
	timestamp TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS "wallet" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL,
	name TEXT
);
CREATE TABLE IF NOT EXISTS "token" (
	id SERIAL PRIMARY KEY,
	address TEXT UNIQUE NOT NULL
);

DO $$ 
    BEGIN
        BEGIN
            ALTER TABLE "token" ADD COLUMN symbol TEXT DEFAULT '';
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column symbol already exists in token.';
        END;
    END;
$$;

CREATE TABLE IF NOT EXISTS "reserve" (
	id SERIAL PRIMARY KEY,
	address TEXT NOT NULL,
	reserve_id TEXT DEFAULT '',
	reserve_type INTEGER DEFAULT 0,
	rebate_wallet TEXT DEFAULT '', 
	block_number INTEGER DEFAULT 0,
	CONSTRAINT reserve_pk UNIQUE (address, reserve_id, block_number)
);



DO $$ 
    BEGIN
        BEGIN
            ALTER TABLE "reserve" ADD COLUMN name TEXT DEFAULT '';
        EXCEPTION
            WHEN duplicate_column THEN RAISE NOTICE 'column name already exists in reserve.';
        END;
    END;
$$;

CREATE TABLE IF NOT EXISTS "` + TradeLogsTableName + `" (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMPTZ,
	block_number INTEGER,
	tx_hash TEXT,
	eth_amount FLOAT(32),
	original_eth_amount FLOAT(32),
	user_address_id BIGINT NOT NULL REFERENCES users,
	src_address_id BIGINT NOT NULL REFERENCES token,
	dst_address_id BIGINT NOT NULL REFERENCES token,
	src_amount FLOAT(32),
	dst_amount FLOAT(32),
	integration_app TEXT,
	ip TEXT,
	country TEXT,
	eth_usd_rate FLOAT(32),
	eth_usd_provider TEXT,
	index INTEGER,
	kyced BOOLEAN,
	is_first_trade BOOLEAN,
	tx_sender	TEXT,
	receiver_address	TEXT,
	gas_used INTEGER,
	gas_price FLOAT(32),
	transaction_fee FLOAT(32),
	version integer,
	CONSTRAINT tradelog_constraint UNIQUE (tx_hash, index)
);

CREATE UNIQUE INDEX IF NOT EXISTS "tradelogs_id_index" ON "` + TradeLogsTableName + `"(id);

ALTER TABLE "` + TradeLogsTableName + `"
	ADD COLUMN IF NOT EXISTS gas_used INTEGER,
	ADD COLUMN IF NOT EXISTS transaction_fee FLOAT(32),
	ADD COLUMN IF NOT EXISTS gas_price FLOAT(32);

CREATE TABLE IF NOT EXISTS "` + BigTradeLogsTableName + `" (
	id SERIAL PRIMARY KEY,
	tradelog_id INTEGER UNIQUE NOT NULL REFERENCES tradelogs (id),
	twitted BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS "trade_timestamp" ON "` + TradeLogsTableName + `"(timestamp);
CREATE INDEX IF NOT EXISTS "trade_user_address" ON "` + TradeLogsTableName + `"(user_address_id);
CREATE INDEX IF NOT EXISTS "trade_src_address" ON "` + TradeLogsTableName + `"(src_address_id);
CREATE INDEX IF NOT EXISTS "trade_dst_address" ON "` + TradeLogsTableName + `"(dst_address_id);
CREATE INDEX IF NOT EXISTS "trade_tx_hash" ON "` + TradeLogsTableName + `"(tx_hash);


CREATE TABLE IF NOT EXISTS "fee" (
	id SERIAL,
	trade_id INTEGER NOT NULL REFERENCES tradelogs,
	reserve_address TEXT NOT NULL,
	wallet_address TEXT default '',
	wallet_fee FLOAT(32) default 0,
	platform_fee FLOAT(32) default 0,
	burn FLOAT(32) default 0,
	rebate FLOAT(32) default 0,
	reward FLOAT(32) default 0,
	version INTEGER default 0
);



-- create_or_update_tradelogs creates or update tradelogs
CREATE OR REPLACE FUNCTION create_or_update_tradelogs(INOUT _id tradelogs.id%TYPE,
												_timestamp tradelogs.timestamp%TYPE,
												_block_number tradelogs.block_number%TYPE,
												_tx_hash tradelogs.tx_hash%TYPE,
												_eth_amount tradelogs.eth_amount%TYPE,
												_original_eth_amount tradelogs.original_eth_amount%TYPE,
												_user_address TEXT,
												_src_address TEXT,
												_dst_address TEXT,
												_src_amount tradelogs.src_amount%TYPE,
												_dst_amount tradelogs.dst_amount%TYPE,
												_integration_app tradelogs.integration_app%TYPE,
												_ip tradelogs.ip%TYPE,
												_country tradelogs.country%TYPE,
												_eth_usd_rate tradelogs.eth_usd_rate%TYPE,
												_eth_usd_provider tradelogs.eth_usd_provider%TYPE,
												_index tradelogs.index%TYPE,
												_kyced tradelogs.kyced%TYPE,
												_is_first_trade tradelogs.is_first_trade%TYPE,
												_tx_sender tradelogs.tx_sender%TYPE,
												_receiver_address tradelogs.receiver_address%TYPE,
												_gas_used tradelogs.gas_used%TYPE,
												_gas_price tradelogs.gas_price%TYPE,
												_transaction_fee tradelogs.transaction_fee%TYPE,
												_version tradelogs.version%TYPE,
												_reserve_addresses TEXT[],
												_platform_wallets TEXT[],
												_platform_fees FLOAT[],
												_burns FLOAT[],
												_rebates FLOAT[],
												_rewards FLOAT[]
												) AS
$$
DECLARE
	_address fee.reserve_address%TYPE;
	_iterator INTEGER := 1;
BEGIN
    IF _id = 0 THEN
		INSERT INTO tradelogs (timestamp, block_number, tx_hash, eth_amount, 
			original_eth_amount, user_address_id, src_address_id, dst_address_id, src_amount, dst_amount,
			integration_app, ip, country, eth_usd_rate, eth_usd_provider, index, kyced, is_first_trade, tx_sender,
			receiver_address, gas_used, gas_price, transaction_fee) 
		VALUES (_timestamp,
			_block_number,
			_tx_hash,
			_eth_amount,
			_original_eth_amount,
			(SELECT id FROM users WHERE address=_user_address),
			(SELECT id FROM token WHERE address=_src_address),
			(SELECT id FROM token WHERE address=_dst_address),
			_src_amount,
			_dst_amount,
			_integration_app,
			_ip,
			_country,
			_eth_usd_rate, 
			_eth_usd_provider,
			_index, 
			_kyced,
			_is_first_trade,
			_tx_sender,
			_receiver_address,
			_gas_used,
			_gas_price,
			_transaction_fee
		) ON CONFLICT (tx_hash, index) DO UPDATE SET 
			timestamp = _timestamp
		 RETURNING id INTO _id;
    END IF;


    IF _id IS NOT NULL THEN
        FOREACH _address IN ARRAY _reserve_addresses
            LOOP
				INSERT INTO "fee"(trade_id, 
					reserve_address, 
					wallet_address, 
					platform_fee, 
					burn, 
					rebate, 
					reward)
				VALUES (_id, _address, 
					_platform_wallets[_iterator],
					_platform_fees[_iterator],
					_burns[_iterator],
					_rebates[_iterator],
					_rewards[_iterator]
				);
				_iterator := _iterator+1;
            END LOOP;
    END IF;

    RETURN;
END;
$$ LANGUAGE PLPGSQL;
`

// DefaultDateFormat ...
const DefaultDateFormat = "2006-01-02 15:04:05"
