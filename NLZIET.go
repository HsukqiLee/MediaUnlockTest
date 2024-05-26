package mediaunlocktest

import (
	"io"
	"net/http"
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

func NLZIET(c http.Client) Result {
	resp, err := GET(c, "https://nlziet.nl/cdn-cgi/trace")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: ErrNetwork}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: ErrNetwork}
	}
	s := string(b)
	i := strings.Index(s, "loc=")
	if i == -1 {
		return Result{Status: StatusUnexpected}
	}
	s = s[i+4:]
	i = strings.Index(s, "\n")
	if i == -1 {
		return Result{Status: StatusUnexpected}
	}
	loc := s[:i]


	if SupportNLZIET(loc) {
		return Result{Status: StatusOK, Region: strings.ToLower(loc)}
	}
	return Result{Status: StatusNo}
}
