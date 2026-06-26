package providers

import (
	"MediaUnlockTest/pkg/core"
	"io"
	"strings"
)

func FXNOW(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://fxnow.fxnetworks.com")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	if strings.Contains(string(b), "is not accessible") {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK}
}
