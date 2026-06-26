package providers

import (
	"MediaUnlockTest/pkg/core"

	http "github.com/bogdanfinn/fhttp"
)

func HoyTV(c core.HttpClient) core.Result {
	return core.CheckGETStatus(c, "https://hoytv-live-stream.hoy.tv/ch78/index-fhd.m3u8", core.ResultMap{
		http.StatusForbidden: {Status: core.StatusNo},
		http.StatusOK:        {Status: core.StatusOK},
	}, core.Result{Status: core.StatusUnexpected})
}
