package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func CW_TV(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://www.cwtv.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	s := string(b)
	if strings.Contains(s, "302 Found") {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK}
}
