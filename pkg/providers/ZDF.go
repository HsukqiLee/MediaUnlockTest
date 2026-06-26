package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func ZDF(c core.HttpClient) core.Result {
	return core.CheckDalvikStatus(c, "https://ssl.zdf.de/geo/de/geo.txt", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
