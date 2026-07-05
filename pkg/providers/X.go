package providers

import (
	"MediaUnlockTest/pkg/core"
	"slices"
	"strings"
)

func SupportX(loc string) bool {
	var X_RESTRICTED_COUNTRY = []string{
		"CN", "IR", "MM", "KP", "RU", "TM",
	}
	return !slices.Contains(X_RESTRICTED_COUNTRY, loc)
}

func X(c core.HttpClient) core.Result {
	loc, err := core.GetCloudflareTraceLoc(c, "https://x.com/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp, err := core.GET(c, "https://x.com/")
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

	if SupportX(loc) {
		return core.Result{Status: core.StatusOK, Region: strings.ToLower(loc)}
	}

	return core.Result{Status: core.StatusNo, Region: strings.ToLower(loc)}
}
