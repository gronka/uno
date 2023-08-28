-- the aim table handles a conversation
DROP TABLE IF EXISTS zinc_searches CASCADE;
CREATE TABLE zinc_searches (
	-- mako-query-id is a murmur3 hash of the string returned by AsQueryString()
	search_id UUID NOT NULL,
	search_result_body TEXT,
	time_fetched BIGINT DEFAULT 0,
	PRIMARY KEY (search_id)
);


DROP TABLE IF EXISTS zinc_product_offers_cache CASCADE;
CREATE TABLE zinc_product_offers_cache (
	search_id UUID NOT NULL,
	search_result_body TEXT,
	time_fetched BIGINT DEFAULT 0,
	PRIMARY KEY (search_id)
);
