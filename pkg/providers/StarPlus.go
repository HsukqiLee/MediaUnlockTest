package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

func SupportStarPlus(loc string) bool {
	var STARPLUS_SUPPORT_COUNTRY = []string{
		"BR", "MX", "AR", "CL", "CO", "PE", "UY", "EC", "PA", "CR", "PY", "BO", "GT", "NI", "DO", "SV", "HN", "VE",
	}
	for _, s := range STARPLUS_SUPPORT_COUNTRY {
		if loc == s {
			return true
		}
	}
	return false
}

func StarPlus(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.starplus.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return core.Result{Status: core.StatusFailed}
	}

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp.StatusCode == 302 && (resp.Header.Get("Location") == "https://www.preview.starplus.com/unavailable" || resp.Header.Get("Location") == "https://www.starplus.com/welcome/unavailable") {
		return core.Result{Status: core.StatusNo}
	}

	if resp.StatusCode == 200 {
		re := regexp.MustCompile(`Region:\s+([A-Za-z]{2})`)
		matches := re.FindStringSubmatch(string(body))
		if len(matches) >= 2 {
			if SupportStarPlus(matches[1]) {
				return core.Result{Status: core.StatusOK, Region: strings.ToLower(matches[1])}
			}
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusUnexpected}
	}

	return core.Result{Status: core.StatusUnexpected}
}
