package providers

import (
	"MediaUnlockTest/pkg/core"
)

func Channel9(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://login.nine.com.au", core.ResultMap{
		403: {Status: core.StatusNo},
		302: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
