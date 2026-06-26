package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func ThreeNow(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://bravo-livestream.fullscreen.nz/index.m3u8")
	if err != nil {
		if core.IsWAFBlockError(err) {
			return core.Result{Status: core.StatusBanned, Err: err}
		}
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
