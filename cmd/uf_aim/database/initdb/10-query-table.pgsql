-- the aim table handles a conversation
DROP TABLE IF EXISTS aims CASCADE;
CREATE TABLE aims (
	surfer_id UUID NOT NULL,
	-- can be shopping, talking
	-- aim_series TEXT [],
	running_tree_name TEXT DEFAULT 'greet_new_user',
	running_branch TEXT DEFAULT 'start',
	previous_running_tree_name TEXT DEFAULT 'none',
	previous_running_branch TEXT DEFAULT 'start',
	challenge_word TEXT DEFAULT '',
	challenge_word_counter INT DEFAULT 0,
	chat_platform TEXT DEFAULT 'sip',
	active_query_id UUID NOT NULL,
	PRIMARY KEY (surfer_id)
);


DROP TABLE IF EXISTS queries CASCADE;
CREATE TABLE queries (
	query_id UUID NOT NULL UNIQUE,
	surfer_id UUID NOT NULL,
	-- I think we don't need this id for anything
	-- zinc_search_id UUID NOT NULL,

	-- details
	current_product_id UUID NOT NULL,
	skipped_product_ids UUID[] NOT NULL,
	current_product_price INT DEFAULT 0,
	absolute_max_price INT DEFAULT 0,
	absolute_min_price INT DEFAULT 0,
	count INT DEFAULT 1,

	-- taxonomy
	noun TEXT DEFAULT '',
	category TEXT DEFAULT '',
	adjectives_sorted TEXT[] DEFAULT '{}',
	anti_adjectives_sorted TEXT[] DEFAULT '{}',

	-- purchase details
	shipping_address_id UUID,
	is_building_address BOOL DEFAULT FALSE,
	
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (query_id)
);
CREATE INDEX idx_queries_surfer_id ON queries(surfer_id);


DROP TABLE IF EXISTS products CASCADE;
CREATE TABLE products (
	product_id UUID NOT NULL UNIQUE,
	query_hash UUID NOT NULL,
	bad_match BOOL DEFAULT false,

	-- taxonomy
	noun TEXT DEFAULT '',
	category TEXT DEFAULT '',
	adjectives_sorted TEXT[] DEFAULT '{}',
	anti_adjectives_sorted TEXT[] DEFAULT '{}',

	-- details
	title TEXT DEFAULT '',
	description TEXT DEFAULT '',
	price INT NOT NULL,
	currency TEXT DEFAULT 'USD',
	review_count INT NOT NULL,
	stars DECIMAL NOT NULL,
	shipping_days INT DEFAULT 0,
	shipping_price INT DEFAULT 0,

	-- scrape_engine is zinc or our in house
	scrape_engine TEXT NOT NULL,
	scrape_engine_product_id TEXT NOT NULL,
	store_name TEXT NOT NULL,
	-- store engine abstracts to amazon, walmart or shopify
	store_engine TEXT NOT NULL,
	store_url TEXT NOT NULL,
	product_url TEXT NOT NULL,
	image_url TEXT NOT NULL,

	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (product_id)
);
CREATE INDEX idx_products_query_hash ON products(query_hash);
