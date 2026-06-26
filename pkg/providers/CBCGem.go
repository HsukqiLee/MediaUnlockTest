package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func CBCGem(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.cbc.ca/g/stats/js/cbc-stats-top.js")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusNo}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, `country":"CA"`) {
		return core.Result{Status: core.StatusOK}
	}
	return core.Result{Status: core.StatusNo}
}
