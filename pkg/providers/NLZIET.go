package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func SupportNLZIET(loc string) bool {
	var NLZIET_SUPPORT_COUNTRY = []string{
		"BE", "BG", "CZ", "DK", "DE", "EE", "IE", "EL", "ES", "FR", "HR", "IT", "CY", "LV", "LT", "LU", "HU", "MT", "NL", "AT", "PL", "PT", "RO", "SI", "SK", "FI", "SE",
	}
	for _, s := range NLZIET_SUPPORT_COUNTRY {
		if loc == s {
			return true
		}
	}
	return false
}

func NLZIET(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://nlziet.nl/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	i := strings.Index(s, "loc=")
	if i == -1 {
		return core.Result{Status: core.StatusUnexpected}
	}
	s = s[i+4:]
	i = strings.Index(s, "\n")
	if i == -1 {
		return core.Result{Status: core.StatusUnexpected}
	}
	loc := s[:i]

	if SupportNLZIET(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}
	return core.Result{Status: core.StatusNo}
}
