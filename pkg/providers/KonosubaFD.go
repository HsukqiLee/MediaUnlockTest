package providers

import (
	"MediaUnlockTest/pkg/core"
)

func KonosubaFD(c core.HttpClient) core.Result {
	resp, err := core.RequestRaw(c, "POST", "https://api.konosubafd.jp/api/masterlist", "", core.H{"User-Agent", "pj0007/212 CFNetwork/1240.0.4 Darwin/20.6.0"})
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		200: {Status: core.StatusOK},
		403: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
