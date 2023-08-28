package uf_aim

import "gitlab.com/textfridayy/uno/uf"

type AimTreeName string

//go:generate go run github.com/tomaspavlic/enumall@latest -type=AimTreeName
const (
	AimGreetNewUser       AimTreeName = "greet_new_user"
	AimContactSupport                 = "contact_support"
	AimSearchItem                     = "search_item"
	AimSearchItemVariant1             = "find_item"
	AimSearchItemVariant2             = "search_product"
	AimSearchItemVariant3             = "shop"
	AimDoSomething                    = "do_something"
	AimAskSomething                   = "ask_something"
	AimDescribeSomething              = "describe_something"
	AimShortcut                       = "shortcut"
	AimShipping                       = "shipping"
	AimBlank                          = "blank"
	AimNone                           = "none"
	AimNotUnderstood                  = "not_understood"

	//AimEmpty                          = ""
	AimTrackOrder    = "track_order"
	AimCancelOrder   = "cancel_order"
	AimReorder       = "reorder"
	AimChangeSize    = "change_size"
	AimChangeAddress = "change_address"

	//AimGreetNewUser      AimTreeName = "GreetNewUser"
	//AimSearchItem                    = "SearchItem"
	//AimDoSomething                   = "DoSomething"
	//AimAskSomething                  = "AskSomething"
	//AimDescribeSomething             = "DescribeSomething"
	//AimShortcut                      = "Shortcut"
	//AimBlank                         = "Blank"
	//AimNone                          = "None"
)

var TreeMap = map[AimTreeName]Tree{
	AimGreetNewUser: {
		Name: AimGreetNewUser,
		Branches: []Branch{
			BranchStart,
			BranchGreetNewUserWhatDoYouNeed,
		},
		ShortCircuit: false,
	},
	AimSearchItem: {
		Name: AimSearchItem,
		Branches: []Branch{
			BranchStart,
			BranchSearchRefine,
			BranchSearchSelectShippingAddress,
			BranchSearchSelectShippingMethod,
		},
		ShortCircuit: false,
	},
	AimBlank: {
		Name:         AimBlank,
		Branches:     []Branch{},
		ShortCircuit: false,
	},
	AimNone: {
		Name:         AimNone,
		Branches:     []Branch{},
		ShortCircuit: false,
	},

	AimDoSomething:       StubbedTree,
	AimAskSomething:      StubbedTree,
	AimDescribeSomething: StubbedTree,
}

func (aimInfo *AimInfo) RunningTree() Tree {
	return TreeMap[aimInfo.RunningTreeName]
}

func (aimInfo *AimInfo) IsRunningTreeOutOfBranches() bool {
	return aimInfo.RunningBranch == BranchDone
}

func (ac *AimTreeName) Priority(branch Branch) int {

	uf.Trace("777777777777777777777777777777")
	uf.Trace(string(*ac))

	switch *ac {
	case AimBlank:
		fallthrough
	//case AimEmpty:
	//fallthrough
	case AimNone:
		return 0

	case AimShortcut:
		fallthrough
	case AimContactSupport:
		return 10
	case AimNotUnderstood:
		return 5

	case AimSearchItem:
		fallthrough
	case AimSearchItemVariant1:
		fallthrough
	case AimSearchItemVariant2:
		fallthrough
	case AimSearchItemVariant3:
		if branch == BranchSearchSelectShippingAddress ||
			branch == BranchSearchInputName ||
			branch == BranchSearchInputPostal ||
			branch == BranchSearchInputAddressLine1 ||
			branch == BranchSearchInputAddressLine2 {
			return 99999
		}
		return 50

	case AimGreetNewUser:
		return 900

	default:
		uf.Error("missing priority for tree: " + string(*ac))
		return 1

	}
}
func (previous *AimTreeName) IsResumableAfter(current AimTreeName) bool {
	switch *previous {
	case AimSearchItem:
		fallthrough
	case AimSearchItemVariant1:
		fallthrough
	case AimSearchItemVariant2:
		fallthrough
	case AimSearchItemVariant3:
		if current.IsSearchItemVariant() {
			return false
		} else {
			return true
		}

	default:
		return true
	}
}

func (ac *AimTreeName) IsSearchItemVariant() bool {
	switch *ac {
	case AimSearchItem:
		fallthrough
	case AimSearchItemVariant1:
		fallthrough
	case AimSearchItemVariant2:
		fallthrough
	case AimSearchItemVariant3:
		return true
	default:
		return false
	}
}

func (ac *AimTreeName) IsTreeNameValid() bool {
	if ac.Priority(BranchStart) == 1 {
		return false
	}
	return true
}

type Tree struct {
	Name         AimTreeName
	Branches     []Branch
	ShortCircuit bool
}
