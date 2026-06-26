package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func KayoSports(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://kayosports.com.au/", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected},
		core.H{"Accept", "*/*"},
		core.H{"Accept-Language", "en-US,en;q=0.9"},
		core.H{"Origin", "https://kayosports.com.au"},
		core.H{"Referer", "https://kayosports.com.au/"},
	)
}
