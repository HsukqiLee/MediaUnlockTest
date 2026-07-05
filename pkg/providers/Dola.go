package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"regexp"
	"strings"
)

var dolaRegionRegex = regexp.MustCompile(`(?i)"inner_region"\s*:\s*"([A-Z]{2})"`)

func Dola(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.dola.com/chat/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusBanned}
	}

	b, err := io.ReadAll(resp.Body)
	if err == nil {
		matches := dolaRegionRegex.FindStringSubmatch(string(b))
		if len(matches) > 1 {
			return core.Result{Status: core.StatusOK, Region: strings.ToLower(matches[1])}
		}
	}

	return core.Result{Status: core.StatusUnexpected}
}
