package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func Popcornflix(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://popcornflix-prod.cloud.seachange.com/cms/popcornflix/clientconfiguration/versions/2")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
