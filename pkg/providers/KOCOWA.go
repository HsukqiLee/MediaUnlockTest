package providers

import (
	"MediaUnlockTest/pkg/core"
)

func KOCOWA(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://www.kocowa.com/", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
