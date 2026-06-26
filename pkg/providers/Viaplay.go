package providers

import (
	"MediaUnlockTest/pkg/core"
	"regexp"
	"strings"
)

func extractViaplayRegion1(url string) string {
	re := regexp.MustCompile(`/([a-z]{2})/`)
	matches := re.FindStringSubmatch(string(url))
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractViaplayRegion2(url string) string {
	re := regexp.MustCompile(`viaplay.([a-z]{2})`)
	matches := re.FindStringSubmatch(string(url))
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func Viaplay(c core.HttpClient) core.Result {
	resp1, err := core.GET(c, "https://viaplay.pl")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()

	if resp1.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	if resp1.StatusCode == 302 && resp1.Header.Get("Location") == "/region-blocked" {
		return core.Result{Status: core.StatusNo}
	}

	if resp1.StatusCode == 302 && resp1.Header.Get("Location") == "https://viaplay.pl/pl-pl/" {
		resp2, err := core.GET(c, "https://viaplay.com/")
		if err != nil {
			return core.Result{Status: core.StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()
		if resp2.StatusCode == 404 {
			return core.Result{Status: core.StatusNo}
		}
		if resp2.StatusCode == 302 {
			if region := extractViaplayRegion1(resp2.Header.Get("Location")); region != "" {
				return core.Result{Status: core.StatusOK, Region: strings.ToLower(region)}
			}
			if region := extractViaplayRegion2(resp2.Header.Get("Location")); region != "" {
				return core.Result{Status: core.StatusOK, Region: region}
			}
		}
	}
	return core.Result{Status: core.StatusUnexpected}
}
