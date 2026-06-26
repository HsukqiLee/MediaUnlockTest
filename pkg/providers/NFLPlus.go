package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func NFLPlus(c core.HttpClient) core.Result {
	url := "https://www.nfl.com/plus/"
	resp, err := core.GET(c, url)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "nflgamepass") {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK}
}
