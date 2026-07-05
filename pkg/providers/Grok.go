package providers

import (
	"MediaUnlockTest/pkg/core"
	"slices"
	"strings"
)

func SupportGrok(loc string) bool {
	var GROK_RESTRICTED_COUNTRY = []string{
		"CN", "RU", "IR", "KP", "CU", "SY",
	}
	return !slices.Contains(GROK_RESTRICTED_COUNTRY, loc)
}

func Grok(c core.HttpClient) core.Result {
	loc, err := core.GetCloudflareTraceLoc(c, "https://grok.com/cdn-cgi/trace")
	if err != nil {
		return core.Result{Status: core.StatusErr, Err: err}
	}

	resp, err := core.GET(c, "https://grok.com/")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	status := core.StatusUnexpected
	if resp.StatusCode == 200 {
		status = core.StatusOK
	} else if resp.StatusCode == 403 {
		if SupportGrok(loc) {
			status = core.StatusBanned
		} else {
			status = core.StatusNo
		}
	}

	res := core.Result{Status: status}
	if res.Status == core.StatusOK || res.Status == core.StatusNo {
		res.Region = strings.ToLower(loc)
	}
	return res
}
