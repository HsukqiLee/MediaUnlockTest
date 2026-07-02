package providers

import (
	"MediaUnlockTest/pkg/core"
)

func JioCinema(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://content-jiovoot.voot.com/psapi/", core.ResultMap{
		200: {Status: core.StatusOK},
		474: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
