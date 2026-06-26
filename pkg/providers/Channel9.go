package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func Channel9(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://login.nine.com.au", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusFound:     {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
