package providers

import (
	"MediaUnlockTest/pkg/core"
	"slices"
	"strings"
)

func SupportPerplexity(loc string) bool {
	var PERPLEXITY_RESTRICTED_COUNTRY = []string{
		"CN", "RU", "IR", "KP", "CU", "SY",
	}
	return !slices.Contains(PERPLEXITY_RESTRICTED_COUNTRY, loc)
}

func Perplexity(c core.HttpClient) core.Result {
	loc, err := core.GetCloudflareTraceLoc(c, "https://www.perplexity.ai/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp, err := core.GET(c, "https://www.perplexity.ai/")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		if SupportPerplexity(loc) {
			return core.Result{Status: core.StatusBanned, Region: strings.ToLower(loc)}
		}
		return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
	}

	if SupportPerplexity(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}

	return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
}
