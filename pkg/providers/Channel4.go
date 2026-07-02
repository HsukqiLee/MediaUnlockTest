package providers

import (
	"MediaUnlockTest/pkg/core"
)

func Channel4(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://www.channel4.com/simulcast/channels/C4", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
