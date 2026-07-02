package providers

import (
	"MediaUnlockTest/pkg/core"
)

func SevenPlus(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://7plus.com.au/", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
