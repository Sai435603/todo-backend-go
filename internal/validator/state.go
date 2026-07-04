package validator

import "crypto/subtle"

func OAuthState(expectedState, providedState string) bool {
	if expectedState == "" || providedState == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expectedState), []byte(providedState)) == 1
}
