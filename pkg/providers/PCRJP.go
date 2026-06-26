package providers

import (
	"MediaUnlockTest/pkg/core"
)

// Princess Connect Re:Dive Japan
func PCRJP(c core.HttpClient) core.Result {
	resp, err := core.GET_Dalvik(c, "https://api-priconne-redive.cygames.jp/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return core.Result{Status: core.StatusOK}
	case 403:
		return core.Result{Status: core.StatusNo}
	}
	return core.Result{Status: core.StatusUnexpected}
}
