package mediaunlocktest

import (
	"io"
	"net/http"
	"regexp"
)

func SupportSonyLiv(loc string) bool {
	var SONYLIV_SUPPORT_COUNTRY = []string{
		"AE", "AF", "AT", "AU", "BD", "BE", "BH", "BT", "CA", "CH", "CN", "DE", "DK", "ES", "FI",
		"FR", "GB", "GR", "HK", "ID", "IE", "IN", "IT", "KW", "LK", "MO", "MV", "MY", "NL", "NO",
		"NP", "NZ", "OM", "PH", "PK", "PL", "PT", "QA", "SA", "SE", "SG", "TH", "TW", "US",
	}
	for _, s := range SONYLIV_SUPPORT_COUNTRY {
		if loc == s {
			return true
		}
	}
	return false
}

func extractSonyLivCountryCode(text string) string {
	re := regexp.MustCompile(`country_code:"([A-Z]{2})"`)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func SonyLiv(c http.Client) Result {
	resp, err := GET(c, "https://www.sonyliv.com/signin")
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if err != nil {
		return Result{Status: StatusNetworkErr, Err: err}
	}
	region := extractSonyLivCountryCode(bodyString)
	if region != "" && SupportSonyLiv(region) {
		return Result{Status: StatusOK, Region: region}
	}
	return Result{Status: StatusNo}
}
