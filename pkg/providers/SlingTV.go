package providers

import (
	"net/http"
	"MediaUnlockTest/pkg/core"
)

func SlingTV(c http.Client) core.Result {
	resp, err := core.GET(c, "https://www.sling.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusFound:     {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

