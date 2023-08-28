DROP TABLE IF EXISTS addresses CASCADE;
CREATE TABLE addresses (
	address_id UUID NOT NULL UNIQUE,
	surfer_id UUID NOT NULL,
	is_builder BOOL NOT NULL,
	is_validated BOOL DEFAULT FALSE,
	name TEXT DEFAULT '',
	address_line_1 TEXT DEFAULT '',
	address_line_2 TEXT DEFAULT '',
	postal TEXT DEFAULT '',
	zip4 TEXT DEFAULT '',
	city TEXT DEFAULT '',
	state TEXT DEFAULT '',
	country TEXT DEFAULT '',
	phone TEXT DEFAULT '',
	instructions TEXT DEFAULT '',
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (address_id)
);
CREATE INDEX idx_addresses_surfer_id ON addresses(surfer_id);
