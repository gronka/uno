package policy

import (
	"bytes"

	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
)

func AtLeastOnePolicyMustAllow(policyChain []bool) bool {
	for _, policy := range policyChain {
		if policy {
			return true
		}
	}
	return false
}

func AllPoliciesMustAllow(policyChain []bool) bool {
	for _, policy := range policyChain {
		if !policy {
			return false
		}
	}
	return true
}

func PolicyPublic() bool {
	//TODO: IP restrictions here?
	return true
}

func PolicyNoOne() bool {
	return false
}

func PolicyToddIsSuperAdmin(toddId *gocql.UUID) bool {
	if *toddId == NinesBytes16 {
		return true
	}
	return false
}

func PolicyToddIsLoggedIn(toddId *gocql.UUID) bool {
	return !bytes.Equal(toddId.Bytes(), ZerosBytes)
}

func PolicyToddIsTester(toddId *gocql.UUID) bool {
	if bytes.Equal(toddId.Bytes(), TwosBytes) {
		return true
	}
	return false
}

func PolicyToddIsSurfer(toddId, targetAccountId *gocql.UUID) bool {
	if toddId == targetAccountId {
		return true
	}
	return false
}

func PolicyCorrectSendBlueSecret(gibs Gibs, conf Config) bool {
	if gibs.SendBlueSecret == conf.SendBlueSecret {
		return true
	}
	return false
}

func PolicyUfKeyIsValid(gibs *uf.Gibs, conf Config) bool {
	if gibs.UfKey == conf.UfKey {
		return true
	}
	return false
}
