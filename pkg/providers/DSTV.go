package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func DSTV(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://now.dstv.com/", core.ResultMap{
		http.StatusUnavailableForLegalReasons: {Status: core.StatusNo},
		http.StatusOK:                         {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
