package providers

import (
	"MediaUnlockTest/pkg/core"
	"net/http"
)

func RakutenMagazine(c http.Client) core.Result {
	resp, err := core.GET(c, "https://data-cloudauthoring.magazine.rakuten.co.jp/rem_repository/////////.key")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		http.StatusNotFound:  {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}

