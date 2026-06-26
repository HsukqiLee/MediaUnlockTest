package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func SupportGPT(loc string) bool {
	var GPT_SUPPORT_COUNTRY = []string{
		"AL", "DZ", "AD", "AO", "AG", "AR", "AM", "AU", "AT", "AZ", "BS", "BD", "BB", "BE", "BZ", "BJ", "BT", "BA", "BW", "BR", "BG", "BF", "CV", "CA", "CL", "CO", "KM", "CR", "HR", "CY", "DK", "DJ", "DM", "DO", "EC", "SV", "EE", "FJ", "FI", "FR", "GA", "GM", "GE", "DE", "GH", "GR", "GD", "GT", "GN", "GW", "GY", "HT", "HN", "HU", "IS", "IN", "ID", "IQ", "IE", "IL", "IT", "JM", "JP", "JO", "KZ", "KE", "KI", "KW", "KG", "LV", "LB", "LS", "LR", "LI", "LT", "LU", "MG", "MW", "MY", "MV", "ML", "MT", "MH", "MR", "MU", "MX", "MC", "MN", "ME", "MA", "MZ", "MM", "NA", "NR", "NP", "NL", "NZ", "NI", "NE", "NG", "MK", "NO", "OM", "PK", "PW", "PA", "PG", "PE", "PH", "PL", "PT", "QA", "RO", "RW", "KN", "LC", "VC", "WS", "SM", "ST", "SN", "RS", "SC", "SL", "SG", "SK", "SI", "SB", "ZA", "ES", "LK", "SR", "SE", "CH", "TH", "TG", "TO", "TT", "TN", "TR", "TV", "UG", "AE", "US", "UY", "VU", "ZM", "BO", "BN", "CG", "CZ", "VA", "FM", "MD", "PS", "KR", "TW", "TZ", "TL", "GB", "AQ",
	}
	for _, s := range GPT_SUPPORT_COUNTRY {
		if loc == s {
			return true
		}
	}
	return false
}

func ChatGPT(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://chatgpt.com/cdn-cgi/trace")
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

	resp, err = core.GET(c, "https://chatgpt.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if strings.Contains(string(b), "VPN") {
		return core.Result{Status: core.StatusBanned, Info: "VPN Blocked"}
	}
	if resp.StatusCode == 429 {
		return core.Result{Status: core.StatusRestricted, Region: strings.ToLower(loc), Info: "429 Rate limit"}
	}
	if loc == "T1" {
		return core.Result{Status: core.StatusOK, Region: "tor"}
	}
	if SupportGPT(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}
	return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
}
