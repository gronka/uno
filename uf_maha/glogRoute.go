package uf_maha

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"gitlab.com/textfridayy/uno/uf"
)

func HGlogCreate(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := uf.GlogStruct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))
	uf.Debug(in)

	_, err := pile.Pool.Exec(ctx, `INSERT INTO glogs (
	glog_id,
	todd_id,
	uf_id,
	chat_id,
	surfer_id,
	uf_id,
	code,
	msg,
	service,
	time) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9
)`,
		uf.RandomUUID(),
		in.ToddId,
		in.UfId,
		in.ChatId,
		in.SurferId,
		in.Code,
		in.Msg,
		in.Service,
		in.Time)

	if err != nil {
		msg := fmt.Sprintf("GlogCreate failed: %v", err)
		uf.LoggingError(msg)
	}

	rb.BuildResponse(ctx, &out)
}

func HGlogsGetByUfId(ctx *fasthttp.RequestCtx, pile uf.Pile) {
	in, out := uf.GlogStruct{}, uf.EmptyStruct{}
	gibs, rb := uf.InitGibs(ctx, &pile, &in)
	rb.ExitIfPolicyFails(uf.PolicyUfKeyIsValid(gibs))
	uf.Debug(in)

	rows, err := pile.Pool.Query(ctx,
		`SELECT 
		glog_id,
		todd_id,
		uf_id,
		chat_id,
		surfer_id,
		code,
		level,
		msg,
		service,
		time
	FROM glogs WHERE uf_id = $1`)
	defer rows.Close()
	uf.FlashError(err)

	glogs := make([]uf.GlogStruct, 0)
	for rows.Next() {
		glog := uf.GlogStruct{}
		rows.Scan(
			&glog.GlogId,
			&glog.ToddId,
			&glog.UfId,
			&glog.ChatId,
			&glog.SurferId,
			&glog.Code,
			&glog.Level,
			&glog.Msg,
			&glog.Service,
			&glog.Time,
		)
		glogs = append(glogs, glog)
	}

	rb.BuildResponse(ctx, &out)
}
