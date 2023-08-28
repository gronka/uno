package ut

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
)

type ChatCollection struct {
	Collection []ChatPg
}

type ChatPg struct {
	ChatId         gocql.UUID
	SurferId       gocql.UUID
	Teaser         string
	TeaserSurferId gocql.UUID
	Msgs           []MsgPg `json:"omitempty"`
}

type MsgPg struct {
	ChatId         gocql.UUID
	MsgId          gocql.UUID
	SenderSurferId gocql.UUID
	Recipient      string
	Content        string
	MediaUrl       string
	ChatPlatform   ChatPlatform

	// Code will help us with multilingual support
	Code        string
	TimeCreated int64
	TimeUpdated int64
}

type MsgsPg struct {
	Msgs []MsgPg
}

type MsgsCollection struct {
	Collection []MsgPg
}

type ChatPlatform string

const (
	ChatPlatformLoop ChatPlatform = "loop"
	ChatPlatformSip               = "sip"
	ChatPlatformWeb               = "web"
)

func LChatGetBySurferId(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (ChatPg, error) {
	uf.Trace("LChatGetBySurferId")
	chat := ChatPg{}

	err := gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
	chat_id,
	surfer_id,
	teaser,
	teaser_surfer_id
	FROM chats WHERE surfer_id=$1`,
		surferId,
	).Scan(
		&chat.ChatId,
		&chat.SurferId,
		&chat.Teaser,
		&chat.TeaserSurferId,
	)
	uf.FlashError(err)

	return chat, err
}

func LChatGetById(
	gibs *uf.Gibs,
	chatId gocql.UUID,
) (ChatPg, error) {
	uf.Trace("LChatGetById")
	chat := ChatPg{}

	err := gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT
	chat_id,
	surfer_id,
	teaser,
	teaser_surfer_id
	FROM chats WHERE chat_id=$1`,
		chatId,
	).Scan(
		&chat.ChatId,
		&chat.SurferId,
		&chat.Teaser,
		&chat.TeaserSurferId,
	)
	uf.FlashError(err)

	return chat, err
}

func (msg *MsgPg) LChatSaveNewMsg(gibs *uf.Gibs) error {
	uf.Trace("LChatSaveNewMsg")
	msg.MsgId, _ = gocql.RandomUUID()
	now := uf.NowStamp()

	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO messages (
		msg_id,
		chat_id,
		sender_surfer_id,
		content,
		code,
		chat_platform,
		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		msg.MsgId,
		msg.ChatId,
		msg.SenderSurferId,
		msg.Content,
		msg.Code,
		msg.ChatPlatform,
		now,
		now,
	)
	uf.FlashError(err)

	return err
}

func LChatGetMsgs(
	gibs *uf.Gibs,
	chatId gocql.UUID,
) ([]MsgPg, error) {
	uf.Trace("LChatGetMsgs")
	rows, err := gibs.Pile.Pool.Query(gibs.Ctx,
		`SELECT chat_id, 
		msg_id, 
		sender_surfer_id,
		content,
		code,
		chat_platform,
		time_created,
		time_updated
	FROM messages WHERE chat_id=$1 
	ORDER BY time_created DESC
	LIMIT 100`, chatId)
	defer rows.Close()
	uf.FlashError(err)

	msgsPg := make([]MsgPg, 0)
	for rows.Next() {
		msg := MsgPg{}
		rows.Scan(
			&msg.ChatId,
			&msg.MsgId,
			&msg.SenderSurferId,
			&msg.Content,
			&msg.Code,
			&msg.ChatPlatform,
			&msg.TimeCreated,
			&msg.TimeUpdated,
		)
		msgsPg = append(msgsPg, msg)
	}

	return msgsPg, err
}

func LChatsGetList(gibs *uf.Gibs) ([]ChatPg, error) {
	uf.Trace("LChatGetList")
	rows, err := gibs.Pile.Pool.Query(gibs.Ctx,
		`SELECT 
	chat_id,
	surfer_id,
	teaser,
	teaser_surfer_id
	FROM chats`)
	defer rows.Close()
	uf.FlashError(err)

	chats := make([]ChatPg, 0)
	for rows.Next() {
		chat := ChatPg{}
		rows.Scan(
			&chat.ChatId,
			&chat.SurferId,
			&chat.Teaser,
			&chat.TeaserSurferId,
		)
		chats = append(chats, chat)
	}

	return chats, err
}

func (chat *ChatPg) LUpsert(gibs *uf.Gibs) error {
	uf.Trace("upsert chat")
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO chats (
		chat_id,
		surfer_id,
		teaser,
		teaser_surfer_id,
		time_created,
		time_updated
	) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (chat_id) DO UPDATE SET
		teaser = EXCLUDED.teaser,
		teaser_surfer_id = EXCLUDED.teaser_surfer_id,
		time_updated = EXCLUDED.time_updated`,

		chat.ChatId,
		chat.SurferId,
		chat.Teaser,
		chat.TeaserSurferId,
		uf.NowStamp(),
		uf.NowStamp(),
	)
	uf.FlashError(err)

	return err
}
