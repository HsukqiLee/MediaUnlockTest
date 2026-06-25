package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func ThreeNow(c http.Client) core.Result {
	resp, err := core.GET(c, "https://bravo-livestream.fullscreen.nz/index.m3u8")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

