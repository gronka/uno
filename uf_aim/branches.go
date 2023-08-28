package uf_aim

type Branch string

const (
	BranchStart   Branch = "start"
	BranchDone           = "done"
	BranchStubbed        = "stubbed"

	BranchSearchRefine = "search_refine"

	BranchSearchSelectShippingAddress = "search_select_shipping_address"
	BranchSearchInputName             = "search_input_name"
	BranchSearchInputPostal           = "search_input_postal"
	BranchSearchInputAddressLine1     = "search_input_address_line_1"
	BranchSearchInputAddressLine2     = "search_input_address_line_2"
	BranchSearchSelectPayment         = "search_select_payment"
	BranchSearchChallengeWord         = "search_challenge_word"

	BranchSearchSelectShippingMethod = "search_select_shipping_method"

	BranchSearchPurchase = "search_purchase"

	BranchGreetNewUserWhatDoYouNeed = "GreetNewUserWhatDoYouNeed"
)

var StubbedTree Tree = Tree{
	Branches:     []Branch{BranchStubbed},
	ShortCircuit: true,
}

func (aimInfo *AimInfo) IsLockedBranch() bool {
	switch aimInfo.RunningBranch {
	case BranchSearchInputName:
		fallthrough
	case BranchSearchInputPostal:
		fallthrough
	case BranchSearchInputAddressLine1:
		fallthrough
	case BranchSearchInputAddressLine2:
		fallthrough
	case BranchSearchSelectPayment:
		return true
	}
	return false
}
