package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
	"strings"
)

func PeacockTV(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.peacocktv.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if strings.Contains(resp.Header.Get("location"), "unavailable") {
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusOK}
}

