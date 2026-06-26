package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func SupportClaude(loc string) bool {
	var CLAUDE_SUPPORT_COUNTRY = []string{
		"AL", "DZ", "AD", "AO", "AG", "AR", "AM", "AU", "AT", "AZ", "BS", "BH", "BD", "BB", "BE", "BZ",
		"BJ", "BT", "BO", "BA", "BW", "BR", "BN", "BG", "BF", "BI", "CV", "KH", "CM", "CA", "TD", "CL",
		"CO", "KM", "CG", "CR", "CI", "HR", "CY", "CZ", "DK", "DJ", "DM", "DO", "EC", "EG", "SV", "GQ",
		"EE", "SZ", "FJ", "FI", "FR", "GA", "GM", "GE", "DE", "GH", "GR", "GD", "GT", "GN", "GW", "GY",
		"HT", "HN", "HU", "IS", "IN", "ID", "IQ", "IE", "IL", "IT", "JM", "JP", "JO", "KZ", "KE", "KI",
		"KW", "KG", "LA", "LV", "LB", "LS", "LR", "LI", "LT", "LU", "MG", "MW", "MY", "MV", "MT", "MH",
		"MR", "MU", "MX", "FM", "MD", "MC", "MN", "ME", "MA", "MZ", "NA", "NR", "NP", "NL", "NZ", "NE",
		"NG", "MK", "NO", "OM", "PK", "PW", "PS", "PA", "PG", "PY", "PE", "PH", "PL", "PT", "QA", "RO",
		"RW", "KN", "LC", "VC", "WS", "SM", "ST", "SA", "SN", "RS", "SC", "SL", "SG", "SK", "SI", "SB",
		"ZA", "KR", "ES", "LK", "SR", "SE", "CH", "TW", "TJ", "TZ", "TH", "TL", "TG", "TO", "TT", "TN",
		"TR", "TM", "TV", "UG", "UA", "AE", "GB", "US", "UY", "UZ", "VU", "VA", "VN", "ZM", "ZW",
	}
	for _, s := range CLAUDE_SUPPORT_COUNTRY {
		if loc == s {
			return true
		}
	}
	return false
}

func Claude(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://claude.ai/cdn-cgi/trace")
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
	if loc == "T1" {
		return core.Result{Status: core.StatusOK, Region: "tor"}
	}
	if SupportClaude(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}
	return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
}
