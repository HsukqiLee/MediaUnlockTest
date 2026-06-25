package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func MoviStarPlus(c http.Client) core.Result {
	resp, err := core.GET(c, "https://contratar.movistarplus.es/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	
	if resp.StatusCode == 403 {
		return core.Result{Status: core.StatusNo}
	}
	
	if resp.StatusCode == 200  {
		return core.Result{Status: core.StatusOK}
	}

	return core.Result{Status: core.StatusUnexpected}
}
