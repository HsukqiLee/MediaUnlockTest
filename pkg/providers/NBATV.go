package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func NBA_TV(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.nba.com/watch/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "Service is not available in your region") {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK}
}
