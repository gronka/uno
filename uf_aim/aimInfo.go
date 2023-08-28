package uf_aim

import (
	"strings"

	"gitlab.com/textfridayy/uno/uf"
)

func (aimInfo *AimInfo) String() string {
	return "SurferId: " + aimInfo.Surfer.SurferId.String() +
		", RunningTree: " + string(aimInfo.RunningTreeName) +
		", PreviousRunningTree: " + string(aimInfo.PreviousRunningTreeName)
}

func (aimInfo *AimInfo) ReportError(gibs *uf.Gibs, details string) {
	uf.Error(aimInfo.String() + gibs.String() + details)
}

func (aimInfo *AimInfo) interpretUserMessage(gibs *uf.Gibs) {
	uf.Trace("interpreting aim")
	aimInfo.MsgIn.Content = strings.TrimSpace(aimInfo.MsgIn.Content)
	if aimInfo.IsShortcut() {
		aimInfo.MakoResponse.Function = AimNone
		//NOTE: shortcuts do not need to be processed by GPT
	} else {
		var err error

		working := true
		counter := 0
		for working {
			aimInfo.MakoResponse, err = InterpretAim(gibs, aimInfo.MsgIn.Content)
			if err != nil {
				aimInfo.ReportError(gibs, "interpret aim failed")
			} else {
				if aimInfo.MakoResponse.Function.IsTreeNameValid() {
					// if tree is defined, then we are done
					working = false
				}
			}
			counter++
			if counter > 3 {
				aimInfo.RunningTreeName = AimNotUnderstood
				working = false
			}
		}
	}
}

func (aimInfo *AimInfo) determineRunningTree(gibs *uf.Gibs) {
	mr := aimInfo.MakoResponse
	if aimInfo.RunningTreeName == AimNone {
		aimInfo.RunningTreeName = AimSearchItem
	}

	uf.Trace("9999999999999999999999999999999999")
	uf.Trace(aimInfo.RunningTreeName)
	uf.Trace(aimInfo.PreviousRunningTreeName)
	uf.Trace(mr.Function)
	if aimInfo.IsShortcut() {
		return
	}

	if aimInfo.RunningBranch == BranchDone {
		if mr.Function.Priority(BranchStart) > aimInfo.PreviousRunningTreeName.Priority(aimInfo.PreviousRunningBranch) {
			aimInfo.setRunningTreeToMr(mr)
		} else {
			uf.Trace("load previous tree")
			aimInfo.loadPreviousRunningTree()
		}

	} else {
		// check if the new command can override the current branch
		if mr.Function.Priority(BranchStart) > aimInfo.RunningTreeName.Priority(aimInfo.RunningBranch) {
			uf.Trace("mr has higher priority")
			aimInfo.savePreviousRunningTree()
			aimInfo.setRunningTreeToMr(mr)
		}
	}

}

func (aimInfo *AimInfo) setRunningTreeToMr(mr *MakoResponse) {
	aimInfo.RunningTreeName = mr.Function
	aimInfo.RunningBranch = BranchStart
	mr.Function = AimNone
}

func (aimInfo *AimInfo) savePreviousRunningTree() {
	aimInfo.PreviousRunningTreeName = aimInfo.RunningTreeName
	aimInfo.PreviousRunningBranch = aimInfo.RunningBranch
}

func (aimInfo *AimInfo) loadPreviousRunningTree() {
	// if we start a new search we do not want to load previous tree
	if aimInfo.PreviousRunningTreeName.IsResumableAfter(aimInfo.RunningTreeName) {
		aimInfo.RunningTreeName = aimInfo.PreviousRunningTreeName
		aimInfo.RunningBranch = aimInfo.PreviousRunningBranch
	} else {
		aimInfo.RunningTreeName = AimNone
		aimInfo.RunningBranch = BranchStart
	}

	// either way, reset PreviousRunningTree
	aimInfo.PreviousRunningTreeName = AimNone
	aimInfo.PreviousRunningBranch = BranchStart
}
