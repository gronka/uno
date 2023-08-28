DROP TABLE IF EXISTS carts CASCADE;
CREATE TABLE carts (
	cart_id UUID NOT NULL UNIQUE,
	surfer_id UUID NOT NULL,

	product_ids UUID[],
	counts INT[],
	prices INT[],
	shippings INT[],

	-- purchase details
	shipping_address_id UUID,
	is_builder_address BOOL DEFAULT FALSE,

	-- cheapest_cost INT DEFAULT -1,
	-- cheapest_days INT DEFAULT -1,
	-- fastest_cost INT DEFAULT -1,
	-- fastest_days INT DEFAULT -1,
	-- balanced_cost INT DEFAULT -1,
	-- balanced_days INT DEFAULT -1,
	-- shipping_choice TEXT DEFAULT '',
	
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (cart_id)
);
CREATE INDEX idx_carts_surfer_id ON carts(surfer_id);
