package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

// World Flipper Japan
func WFJP(c core.HttpClient) core.Result {
	return core.CheckDalvikStatus(c, "https://api.worldflipper.jp/", core.ResultMap{
		http.StatusOK:        {Status: core.StatusOK},
		http.StatusForbidden: {Status: core.StatusNo},
	}, core.Result{Status: core.StatusUnexpected})
}
