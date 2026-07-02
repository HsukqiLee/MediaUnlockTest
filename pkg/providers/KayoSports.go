package providers

import (
	"MediaUnlockTest/pkg/core"
)

func KayoSports(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://kayosports.com.au/", core.ResultMap{
		403: {Status: core.StatusNo},
		200: {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected},
		core.H{"Accept", "*/*"},
		core.H{"Accept-Language", "en-US,en;q=0.9"},
		core.H{"Origin", "https://kayosports.com.au"},
		core.H{"Referer", "https://kayosports.com.au/"},
	)
}
