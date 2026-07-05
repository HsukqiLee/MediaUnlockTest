package providers

import (
	"MediaUnlockTest/pkg/core"
	"slices"
	"strings"
)

func SupportMistral(loc string) bool {
	var MISTRAL_RESTRICTED_COUNTRY = []string{
		"RU", "BY", "KP", "IR", "SY", "CU", "CN", "TM",
	}
	return !slices.Contains(MISTRAL_RESTRICTED_COUNTRY, loc)
}

func Mistral(c core.HttpClient) core.Result {
	loc, err := core.GetCloudflareTraceLoc(c, "https://chat.mistral.ai/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp, err := core.GET(c, "https://chat.mistral.ai/")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
	}

	if resp.StatusCode == 200 || resp.StatusCode == 307 || resp.StatusCode == 308 {
		if SupportMistral(loc) {
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
		}
		return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
	}

	return core.Result{Status: core.StatusUnexpected}
}
