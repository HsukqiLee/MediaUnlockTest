package providers

import (
	"MediaUnlockTest/pkg/core"
)

func RakutenMagazine(c core.HttpClient) core.Result {
	resp, err := core.GET(c, "https://data-cloudauthoring.magazine.rakuten.co.jp/rem_repository/////////.key")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		404: {Status: core.StatusOK},
		403: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
