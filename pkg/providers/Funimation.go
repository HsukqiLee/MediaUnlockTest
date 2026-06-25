package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func Funimation(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.crunchyroll.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}
	for _, c := range resp.Cookies() {
		if c.Name == "region" {
			return core.Result{Status: core.StatusOK, Region: c.Value}
		}
	}
	return core.Result{Status: core.StatusFailed}
}

