package providers

import (
	"net/http"
	"MediaUnlockTest/pkg/core"
)

func SHOWTIME(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.paramountpluswithshowtime.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

