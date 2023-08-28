-- the aim table handles a conversation
DROP TABLE IF EXISTS glogs CASCADE;
CREATE TABLE glogs (
	glog_id UUID NOT NULL,
	todd_id UUID,
	uf_id UUID,
	chat_id UUID,
	surfer_id UUID,
	code TEXT,
	level int,
	msg TEXT,
	service TEXT,
	time BIGINT NOT NULL,
	PRIMARY KEY (glog_id)
);
CREATE INDEX idx_glogs_todd_id ON glogs(todd_id);
CREATE INDEX idx_glogs_uf_id ON glogs(uf_id);
CREATE INDEX idx_glogs_code ON glogs(code);
