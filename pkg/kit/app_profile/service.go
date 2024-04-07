package app_profile

import (
	"os"
	"strings"
)

// GetProfileByScopeSuffix retrieves the app_profile by coincidence in the scope suffix
// this is an adaptation of the ScopeUtils.java
func GetProfileByScopeSuffix() string {
	tokens := strings.Split(getScopeValue(), "-")
	return tokens[len(tokens)-1]
}

// IsLocalProfile retrieves information about if the app_profile is local or not
// this is an adaptation of the ScopeUtils.java
func IsLocalProfile() bool {
	return Local == getScopeValue()
}

// IsTestProfile retrieves information about if the app_profile is test or not
// this is an adaptation of the ScopeUtils.java
func IsTestProfile() bool {
	return strings.HasSuffix(getScopeValue(), Test)
}

// IsProdProfile retrieves information about if the app_profile is prod or not
// this is an adaptation of the ScopeUtils.java
func IsProdProfile() bool {
	return strings.HasSuffix(getScopeValue(), Prod)
}

func getScopeValue() string {
	scope := os.Getenv("SCOPE")
	if scope != "" {
		return scope
	}
	return Local
}
