DROP TABLE IF EXISTS stripe_payment_attempts CASCADE;
CREATE TABLE stripe_payment_attempts (
	attempt_id UUID NOT NULL,
	surfer_id UUID NOT NULL,
	order_id UUID,
	payment_intent_id TEXT,
	admin_notes TEXT DEFAULT '',
	credit_applied INT NOT NULL,
	currency TEXT DEFAULT 'USD',
	is_autopay_selected BOOL NOT NULL,
	plastic_id TEXT NOT NULL,
	-- silo can be 'order' or 'plan'
	silo TEXT NOT NULL,
	status INT NOT NULL,
	terms_accepted_version TEXT NOT NULL,
	total_price INT NOT NULL,

	-- total_price_after_credit INT NOT NULL,
	-- bill_id UUID,

	-- time_pulled is used with INSERT SELECT to prevent race conditions
	time_pulled BIGINT DEFAULT 0,
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (attempt_id)
);
CREATE INDEX idx_stripe_payment_attempts_status ON 
	stripe_payment_attempts(status);

