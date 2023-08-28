package uf

import (
	"bytes"

	"github.com/gocql/gocql"
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

func PolicyToddIsSuperAdmin(toddId gocql.UUID) bool {
	Trace("toddId: " + toddId.String())
	if bytes.Equal(toddId.Bytes(), SevensBytes) ||
		bytes.Equal(toddId.Bytes(), NinesBytes) {
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

func PolicyCorrectSendBlueHookSecret(gibs Gibs) bool {
	if gibs.SendBlueHookSecret == gibs.Conf.SendBlueHookSecret {
		return true
	}
	return false
}

func PolicyCorrectLoopHookAuthorization(gibs Gibs) bool {
	if gibs.LoopHookAuthorization == gibs.Conf.LoopHookAuthorization {
		return true
	}
	return false
}

func PolicyCorrectSipHookSecret(gibs Gibs) bool {
	if gibs.SipHookSecret == gibs.Conf.SipHookSecret {
		return true
	}
	return false
}

func PolicyUfKeyIsValid(gibs Gibs) bool {
	if gibs.UfKey == gibs.Conf.UfKey {
		return true
	}
	return false
}
