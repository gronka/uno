DROP TABLE IF EXISTS orders CASCADE;
CREATE TABLE orders (
	order_id UUID NOT NULL UNIQUE,
	surfer_id UUID NOT NULL,
	surfer_subscription_id INT NOT NULL,
	basket_item_ids UUID[] NOT NULL,

	payment_attempt_ids UUID[],
	was_payment_made BOOL DEFAULT FALSE,
	status TEXT NOT NULL,

	title TEXT NOT NULL,
	basket_price INT NOT NULL,
	tax TEXT NOT NULL,
	margin INT NOT NULL,
	credit_to_apply INT NOT NULL,
	currency TEXT DEFAULT 'USD',
	shipping_days INT NOT NULL,
	shipping_price INT NOT NULL,

	shipping_address_id UUID NOT NULL,
	shipping_name TEXT NOT NULL,
	shipping_line_1 TEXT NOT NULL,
	shipping_line_2 TEXT NOT NULL,
	shipping_city TEXT NOT NULL,
	shipping_state TEXT NOT NULL,
	shipping_postal TEXT NOT NULL,

	driver TEXT, -- would be zinc or shopify
	zinc_order_request_id TEXT DEFAULT '',

	time_placed BIGINT DEFAULT 0,
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (order_id)
);
CREATE INDEX idx_orders_surfer_id ON orders(surfer_id);


DROP TABLE IF EXISTS basket_items CASCADE;
CREATE TABLE basket_items (
	basket_item_id UUID NOT NULL UNIQUE,
	order_id UUID NOT NULL,
	product_id UUID NOT NULL,
	title TEXT NOT NULL,
	quantity INT NOT NULL,
	unit_price INT NOT NULL,
	currency TEXT NOT NULL,

	scrape_engine TEXT NOT NULL,
	scrape_engine_product_id TEXT NOT NULL,
	store_url TEXT,
	product_url TEXT,
	image_url TEXT,

	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (basket_item_id)
);


DROP TABLE IF EXISTS payment_attempts CASCADE;
CREATE TABLE payment_attempts (
	payment_attempt_id UUID NOT NULL UNIQUE,
	order_id UUID NOT NULL,

	driver TEXT NOT NULL,
	is_autopay_selected BOOL NOT NULL,
	status TEXT NOT NULL,
	notes TEXT NOT NULL,
	stripe_payment_intent_id TEXT NOT NULL,
	stripe_payment_method_id TEXT NOT NULL,

	price INT NOT NULL,
	currency TEXT NOT NULL,

	terms_accepted_version TEXT NOT NULL,
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (payment_attempt_id)
);
CREATE INDEX idx_payment_attempts_stripe_payment_intent_id ON 
	payment_attempts(stripe_payment_intent_id);
