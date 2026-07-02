package providers

import (
	"MediaUnlockTest/pkg/core"
)

func CoupangPlay(c core.HttpClient) core.Result {
	resp, err := core.GET_Dalvik(c, "https://www.coupangplay.com/")
	if err != nil {
		return core.Result{Status: core.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 && resp.Header.Get("Location") == "https://www.coupangplay.com/not-available" {
		return core.Result{Status: core.StatusNo}
	}

	return core.ResultFromMapping(resp.StatusCode, core.ResultMap{
		200: {Status: core.StatusOK},
		403: {Status: core.StatusBanned},
	}, core.Result{Status: core.StatusUnexpected})
}
