package uf_aim

import (
	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uxt"
)

func (aimInfo *AimInfo) SetReplyIfUserWantsToQuitShipping(gibs *uf.Gibs) bool {
	quit := aimInfo.IsShortcutQuit()
	if quit {
		aimInfo.MsgOut.Content = "Okay, we can do that later. Let's go " +
			"back to searching."
		uxt.AddressDeleteBuilder(gibs, aimInfo.Surfer.SurferId)
		aimInfo.RunningBranch = BranchSearchRefine
		aimInfo.MustWaitUserResponse = true
		return true
	}
	return false
}
