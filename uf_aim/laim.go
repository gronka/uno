package uf_aim

import (
	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
)

func LAimGetBySurferId(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (aimInfo AimInfo, err error) {

	uf.Trace("aim select")
	err = gibs.Pile.Pool.QueryRow(gibs.Ctx, `SELECT 
		surfer_id,
		running_tree_name,
		running_branch,
		previous_running_tree_name,
		previous_running_branch,
		active_query_id
	FROM aims WHERE surfer_id=$1`,
		surferId,
	).Scan(
		&aimInfo.SurferId,
		&aimInfo.RunningTreeName,
		&aimInfo.RunningBranch,
		&aimInfo.PreviousRunningTreeName,
		&aimInfo.PreviousRunningBranch,
		&aimInfo.ActiveQueryId,
	)
	uf.FlashError(err)

	uf.Debug(aimInfo.String())

	return aimInfo, err
}

func LAimCreate(
	gibs *uf.Gibs,
	surferId gocql.UUID,
) (err error) {

	uf.Trace("aim create")
	_, err = gibs.Pile.Pool.Exec(gibs.Ctx, `INSERT INTO aims (
		surfer_id, active_query_id) VALUES ($1, $2)`,
		surferId,
		uf.ZerosUuid,
	)
	uf.FlashError(err)

	//if !token.Insert() {
	//return aim, errors.New("insert aim failed")
	//}

	return err
}

func (aimInfo *AimInfo) LAimUpdate(gibs *uf.Gibs) error {
	_, err := gibs.Pile.Pool.Exec(gibs.Ctx, `UPDATE aims SET
		running_tree_name = $1,
		running_branch = $2,
		previous_running_tree_name = $3,
		previous_running_branch = $4,
		active_query_id = $5,
		challenge_word = $6,
		challenge_word_counter = $7,
		chat_platform = $8
	WHERE surfer_id = $9`,
		aimInfo.RunningTreeName,
		aimInfo.RunningBranch,
		aimInfo.PreviousRunningTreeName,
		aimInfo.PreviousRunningBranch,
		aimInfo.ActiveQueryId,
		aimInfo.ChallengeWord,
		aimInfo.ChallengeWordCounter,
		aimInfo.ChatPlatform,
		aimInfo.Surfer.SurferId,
	)
	uf.FlashError(err)

	//if !token.Update() {
	//return errors.New("update aim failed")
	//}

	return err
}
