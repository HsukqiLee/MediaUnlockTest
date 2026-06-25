package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func PandaTV(c http.Client) core.Result {
	resp, err := core.GET(c, "https://api.pandalive.co.kr/v1/live/play")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusBadRequest: {Status: core.StatusOK},
		http.StatusForbidden:  {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

