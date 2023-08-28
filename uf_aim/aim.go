package uf_aim

import (
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

type AimInfo struct {
	StepCounter int

	Surfer    uf_user.SurferPg
	Recipient string
	Chat      ut.ChatPg
	MsgIn     ut.MsgPg
	MsgOut    ut.MsgPg
	// Iterating is a flag which we can set to True in case Fridayy needs to
	// handle multiple actions before completion.
	Iterating                       bool
	MustResumePreviousRunningBranch bool
	MustWaitUserResponse            bool

	// searching members
	MakoResponse *MakoResponse
	QueryFilters *QueryFilters
	Error        error
	Query        *QueryPg

	Cart *CartPg

	// variables stored in postgresql
	SurferId                gocql.UUID
	RunningTreeName         AimTreeName
	RunningBranch           Branch
	PreviousRunningTreeName AimTreeName
	PreviousRunningBranch   Branch
	ActiveQueryId           gocql.UUID

	ChallengeWord        string
	ChallengeWordCounter int
	ChatPlatform         ut.ChatPlatform
}

type QueryFilters struct {
}

// maxSteps is to avoid infinite loop. It may need tuning
var maxSteps = 5

func (aimInfo *AimInfo) GetMsgIn() string {
	if aimInfo.StepCounter == 0 {
		return aimInfo.MsgIn.Content
	} else {
		return ""
	}
}

func (aimInfo *AimInfo) IsShortcutSansCounter() bool {
	if IsShortcut(aimInfo.MsgIn.Content) {
		return true
	}
	return false
}

func (aimInfo *AimInfo) IsShortcut() bool {
	if aimInfo.StepCounter == 0 {
		if IsShortcut(aimInfo.MsgIn.Content) {
			return true
		}
	}
	return false
}

func (aimInfo *AimInfo) Process(
	gibs *uf.Gibs,
) (ures uf.UfResponse, err error) {

	//TODO: should we always allow 'quit' to reset to search branch?
	aimInfo.StepCounter = 0
	if aimInfo.IsShortcutSansCounter() {
		aimInfo.MakoResponse.Function = AimShortcut
	} else if aimInfo.IsLockedBranch() {
		// do nothing
		aimInfo.MakoResponse.Function = AimShortcut
	} else {
		uf.Trace("Interpret: " + aimInfo.MsgIn.Content)
		aimInfo.interpretUserMessage(gibs)
		uf.Trace(aimInfo.MakoResponse)
	}

	for aimInfo.Iterating {
		// before anything, reset msgout
		aimInfo.MsgOut = ut.MsgPg{
			ChatId:         aimInfo.Chat.ChatId,
			SenderSurferId: uf.FridayyUuid,
			Recipient:      aimInfo.Recipient,
			ChatPlatform:   aimInfo.ChatPlatform,
			Content:        "",
		}

		uf.Trace("iteration: " + strconv.Itoa(aimInfo.StepCounter))
		aimInfo.determineRunningTree(gibs)
		uf.Trace(aimInfo)

		if aimInfo.StepCounter > 0 {
			// space out operations a bit
			time.Sleep(40 * time.Millisecond)
		}

		ures, err = aimInfo.handleNextIteration(gibs)
		if err != nil {
			aimInfo.ReportError(gibs, "handleNextIteration failed: "+err.Error())
		}

		uf.Trace("approaching end of branch")
		if aimInfo.RunningTreeName == AimNone ||
			aimInfo.MustWaitUserResponse ||
			aimInfo.StepCounter == maxSteps {
			aimInfo.Iterating = false
		}
		aimInfo.StepCounter++
	}
	uf.Trace("iterating done")
	return ures, nil
}

func (aimInfo *AimInfo) handleNextIteration(
	gibs *uf.Gibs,
) (ures uf.UfResponse, err error) {
	if aimInfo.MustResumePreviousRunningBranch {
		ures, err = aimInfo.resumePreviousRunningBranch(gibs)
		aimInfo.MustResumePreviousRunningBranch = false
	} else {
		ures, err = aimInfo.processTree(gibs)
	}

	if aimInfo.MsgOut.Content != "" {
		uf.Trace("about to save")
		err = aimInfo.MsgOut.LChatSaveNewMsg(gibs)
		if err != nil {
			uf.Error(err)
			ures.AddError(SaveMessageError)
		}
		ures = uxt.SendReply(gibs, &aimInfo.MsgOut)
	}

	uf.Trace("about to update aim")
	err = aimInfo.LAimUpdate(gibs)
	if err != nil {
		uf.Error(err)
		ures.AddError(AimUpdateError)
	}

	uf.Trace("msgOut: " + aimInfo.MsgOut.Content)
	uf.Trace("ures.Errors:")
	uf.Trace(ures.Errors)

	return ures, err
}
