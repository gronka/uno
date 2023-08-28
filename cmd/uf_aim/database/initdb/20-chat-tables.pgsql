DROP TABLE IF EXISTS chats CASCADE;
CREATE TABLE chats (
	chat_id UUID NOT NULL,
	surfer_id UUID NOT NULL,
	teaser TEXT DEFAULT '',
	teaser_surfer_id UUID NOT NULL,

	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (chat_id)
);
CREATE INDEX idx_chats_surfer_id ON chats(surfer_id);


DROP TABLE IF EXISTS messages CASCADE;
CREATE TABLE messages (
	msg_id UUID NOT NULL UNIQUE,
	chat_id UUID NOT NULL,
	sender_surfer_id UUID NOT NULL,
	content TEXT DEFAULT '',
	code TEXT DEFAULT '',
	media_url TEXT DEFAULT '',
	chat_platform TEXT DEFAULT '',
	
	time_created BIGINT DEFAULT 0,
	time_updated BIGINT DEFAULT 0,
	PRIMARY KEY (msg_id)
);
CREATE INDEX idx_messages_chat_id ON messages(chat_id, time_created);
