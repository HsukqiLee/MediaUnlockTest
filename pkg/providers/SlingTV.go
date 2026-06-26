package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func SlingTV(c core.HttpClient) core.Result {
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
