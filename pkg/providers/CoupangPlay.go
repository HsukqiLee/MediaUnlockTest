package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
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
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusBanned},
	}, core.Result{Status: core.StatusUnexpected})
}
