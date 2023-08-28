package uf_aim

import (
	"gitlab.com/textfridayy/uno/uf"
)

// This function will probably simply be to send a remind to the user of where
// they are in the process of shopping with Fridayy
func (aimInfo *AimInfo) resumePreviousRunningBranch(
	gibs *uf.Gibs,
) (ures uf.UfResponse, err error) {
	switch aimInfo.RunningTreeName {

	case AimSearchItem:
		aimInfo.MsgOut.Content = "Okay! Back to shopping!"

	case AimDoSomething:
	case AimAskSomething:
	case AimDescribeSomething:
	case AimBlank:
	case AimGreetNewUser:

	default:
		uf.Error("Unknown AimSeries case")
		aimInfo.RunningTreeName = AimBlank
	}

	return ures, err
}
