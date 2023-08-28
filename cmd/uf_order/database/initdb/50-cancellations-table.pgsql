DROP TABLE IF EXISTS cancellations CASCADE;
CREATE TABLE cancellations (
	cancellation_id UUID NOT NULL UNIQUE,
	zinc_cancel_request_id TEXT,
	zinc_order_request_id TEXT,
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (cancellation_id)
);
CREATE INDEX idx_cancellations_zinc_order_request_id ON 
	cancellations(zinc_order_request_id);
CREATE INDEX idx_cancellations_zinc_cancel_request_id ON 
	cancellations(zinc_cancel_request_id);
