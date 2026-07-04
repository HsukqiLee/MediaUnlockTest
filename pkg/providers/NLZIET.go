package providers

import (
	"MediaUnlockTest/pkg/core"
	"slices"
	"strings"
)

func SupportNLZIET(loc string) bool {
	var NLZIET_SUPPORT_COUNTRY = []string{
		"BE", "BG", "CZ", "DK", "DE", "EE", "IE", "EL", "ES", "FR", "HR", "IT", "CY", "LV", "LT", "LU", "HU", "MT", "NL", "AT", "PL", "PT", "RO", "SI", "SK", "FI", "SE",
	}
	return slices.Contains(NLZIET_SUPPORT_COUNTRY, loc)
}

func NLZIET(c core.HttpClient) core.Result {
	loc, err := core.GetCloudflareTraceLoc(c, "https://nlziet.nl/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	if SupportNLZIET(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}
	return core.Result{Status: core.StatusNo}
}
